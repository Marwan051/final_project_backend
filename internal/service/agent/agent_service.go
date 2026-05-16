package agent

import "context"

// Request represents an agent message exchange request.
type Request struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
}

// Response represents an agent message exchange response.
type Response struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
}

// Agent defines the agent service contract used by the HTTP layer.
type Agent interface {
	Query(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
