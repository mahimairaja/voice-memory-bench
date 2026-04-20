package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a thin HTTP client that talks to a local sidecar over 127.0.0.1.
// All latency measurements on the read-path are taken by the caller around
// Do* methods — this file deliberately does not try to time the wire itself.
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient constructs a client pointed at the given base URL
// (for example http://127.0.0.1:8765).
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Health pings GET /health.
func (c *Client) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health check failed: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

// Capabilities fetches GET /capabilities.
func (c *Client) Capabilities(ctx context.Context) (*Capabilities, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/capabilities", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var caps Capabilities
	if err := json.NewDecoder(resp.Body).Decode(&caps); err != nil {
		return nil, fmt.Errorf("decode capabilities: %w", err)
	}
	return &caps, nil
}

// AddMessage calls POST /add_message.
func (c *Client) AddMessage(ctx context.Context, req AddMessageRequest) (*WriteResult, error) {
	var out WriteResult
	if err := c.post(ctx, "/add_message", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFact calls POST /add_fact.
func (c *Client) AddFact(ctx context.Context, req AddFactRequest) (*WriteResult, error) {
	var out WriteResult
	if err := c.post(ctx, "/add_fact", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search calls POST /search.
func (c *Client) Search(ctx context.Context, req SearchRequest) (*RetrievalResult, error) {
	var out RetrievalResult
	if err := c.post(ctx, "/search", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Enumerate calls POST /enumerate.
func (c *Client) Enumerate(ctx context.Context, req EnumerateRequest) ([]MemoryItem, error) {
	var out struct {
		Items []MemoryItem `json:"items"`
	}
	if err := c.post(ctx, "/enumerate", req, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

// Reset calls POST /reset.
func (c *Client) Reset(ctx context.Context, req ResetRequest) error {
	return c.post(ctx, "/reset", req, nil)
}

func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s request: %w", path, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return err
	}
	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s response: %w", path, err)
	}
	return nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	var env ErrorEnvelope
	if err := json.Unmarshal(body, &env); err == nil && env.ErrorType != "" {
		if resp.StatusCode == http.StatusUnprocessableEntity && env.ErrorType == "capability_not_supported" {
			return &CapabilityNotSupportedError{
				Provider:   env.Provider,
				Capability: env.Capability,
				Reason:     env.Message,
			}
		}
		return &SidecarError{StatusCode: resp.StatusCode, ErrorType: env.ErrorType, Message: env.Message}
	}
	return &SidecarError{StatusCode: resp.StatusCode, ErrorType: "unknown", Message: string(body)}
}
