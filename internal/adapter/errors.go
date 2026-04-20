package adapter

import "fmt"

// CapabilityNotSupportedError is raised when the sidecar replies 422 with
// error_type=capability_not_supported. The engine treats this as "skip with
// reason" rather than a hard failure.
type CapabilityNotSupportedError struct {
	Provider   string
	Capability string
	Reason     string
}

func (e *CapabilityNotSupportedError) Error() string {
	return fmt.Sprintf("[%s] capability %q not supported: %s", e.Provider, e.Capability, e.Reason)
}

// SidecarError is the generic sidecar failure.
type SidecarError struct {
	StatusCode int
	ErrorType  string
	Message    string
}

func (e *SidecarError) Error() string {
	return fmt.Sprintf("sidecar error %d (%s): %s", e.StatusCode, e.ErrorType, e.Message)
}

// ErrorEnvelope is the JSON body the sidecar returns on 4xx/5xx.
type ErrorEnvelope struct {
	ErrorType  string `json:"error_type"`
	Message    string `json:"message"`
	Provider   string `json:"provider,omitempty"`
	Capability string `json:"capability,omitempty"`
}
