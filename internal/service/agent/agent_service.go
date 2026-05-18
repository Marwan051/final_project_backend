package agent

import "context"

// AgentRequest represents an agent message exchange request.
type AgentRequest struct {
	UserQuery string `json:"user_query"`
	SessionID string `json:"session_id,omitempty"`
}

// AgentResponse represents an agent message exchange response.
type AgentResponse struct {
	Answer    string `json:"answer"`
	SessionID string `json:"session_id"`
}

// Agent defines the agent service contract used by the HTTP layer.
type Agent interface {
	ProcessQuery(ctx context.Context, req AgentRequest) (AgentResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
