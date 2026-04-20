package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/mahimairaja/vbench/internal/adapter"
)

// Process wraps a running sidecar subprocess.
type Process struct {
	cmd    *exec.Cmd
	port   int
	client *adapter.Client
}

// Client returns an HTTP client pointed at the sidecar.
func (p *Process) Client() *adapter.Client { return p.client }

// Port returns the port the sidecar is bound to.
func (p *Process) Port() int { return p.port }

// SpawnOptions controls how a sidecar subprocess is launched.
type SpawnOptions struct {
	// Command is the argv for the sidecar (e.g. ["uv", "run", "vbench-mem0"]).
	Command []string
	// WorkingDir is the directory to run the command in (the sidecar package root).
	WorkingDir string
	// ProviderConfig is serialised to JSON and passed via VBENCH_PROVIDER_CONFIG.
	ProviderConfig map[string]interface{}
	// ExtraEnv is merged onto os.Environ().
	ExtraEnv []string
	// ReadyTimeout bounds how long we wait for /health to succeed.
	ReadyTimeout time.Duration
	// Stdout / Stderr are used for child output; defaults to os.Stdout / os.Stderr.
	Stdout *os.File
	Stderr *os.File
}

// Spawn starts a sidecar subprocess, waits for /health, and returns a Process handle.
func Spawn(ctx context.Context, opts SpawnOptions) (*Process, error) {
	if len(opts.Command) == 0 {
		return nil, fmt.Errorf("sidecar command is empty")
	}
	port, err := freePort()
	if err != nil {
		return nil, fmt.Errorf("allocate sidecar port: %w", err)
	}

	configJSON, err := json.Marshal(opts.ProviderConfig)
	if err != nil {
		return nil, fmt.Errorf("marshal provider config: %w", err)
	}

	cmd := exec.CommandContext(ctx, opts.Command[0], opts.Command[1:]...)
	cmd.Dir = opts.WorkingDir
	cmd.Env = append(os.Environ(),
		"VBENCH_SIDECAR_PORT="+strconv.Itoa(port),
		"VBENCH_PROVIDER_CONFIG="+string(configJSON),
	)
	cmd.Env = append(cmd.Env, opts.ExtraEnv...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start sidecar: %w", err)
	}

	client := adapter.NewClient(fmt.Sprintf("http://127.0.0.1:%d", port))
	readyTimeout := opts.ReadyTimeout
	if readyTimeout == 0 {
		readyTimeout = 30 * time.Second
	}
	if err := waitReady(ctx, client, readyTimeout); err != nil {
		_ = terminate(cmd)
		return nil, fmt.Errorf("sidecar did not become ready within %s: %w", readyTimeout, err)
	}

	return &Process{cmd: cmd, port: port, client: client}, nil
}

// Shutdown terminates the sidecar subprocess and its process group.
func (p *Process) Shutdown() error {
	if p == nil || p.cmd == nil {
		return nil
	}
	return terminate(p.cmd)
}

func terminate(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		if pgid > 0 {
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			_ = cmd.Process.Kill()
		}
		<-done
	}
	return nil
}

func waitReady(ctx context.Context, client *adapter.Client, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("deadline exceeded")
		}
		pingCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		err := client.Health(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
}

// freePort asks the kernel for an ephemeral port by binding to :0, then
// closes the listener and returns the port so the sidecar subprocess can
// bind it moments later. This is a narrow TOCTOU race — another process on
// the same loopback could steal the port between Close() and the sidecar's
// bind — but for local benchmarking on 127.0.0.1 the window is microseconds
// wide and the impact is a retryable startup failure, not data loss. The
// alternative (passing a pre-bound listener into uvicorn via the sidecar) is
// not worth the complexity at this stage.
//
// net.Listen is called without a context because the bind is effectively
// non-blocking on loopback. Listener Close errors on an ephemeral port are
// not actionable and are intentionally ignored.
func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
