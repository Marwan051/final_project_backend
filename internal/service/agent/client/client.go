package client

import (
	"context"
	"fmt"
	"time"

	authpkg "github.com/Marwan051/final_project_backend/internal/server"
	agent_service "github.com/Marwan051/final_project_backend/internal/service/agent"
	pb "github.com/Marwan051/final_project_backend/internal/service/agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type AgentClient struct {
	client         pb.AgentServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewAgentClient(cfg ClientConfig) (*AgentClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	opts = append(opts, cfg.DialOptions...)

	conn, err := grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client for %s: %w", cfg.Address, err)
	}

	timeout := cfg.RequestTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &AgentClient{
		client:         pb.NewAgentServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *AgentClient) Close() error {
	return c.conn.Close()
}

func (c *AgentClient) ProcessQuery(ctx context.Context, req agent_service.AgentRequest) (agent_service.AgentResponse, error) {
	if err := ctx.Err(); err != nil {
		return agent_service.AgentResponse{}, err
	}

	if token, ok := authpkg.GetAuthToken(ctx); ok && token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	resp, err := c.client.ProcessQuery(ctx, &pb.AgentRequest{
		UserQuery: req.UserQuery,
		SessionId: req.SessionID,
	})
	if err != nil {
		return agent_service.AgentResponse{}, fmt.Errorf("grpc query failed: %w", err)
	}

	return agent_service.AgentResponse{
		Answer:    resp.Answer,
		SessionID: resp.SessionId,
	}, nil
}

func (c *AgentClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		return false, err
	}
	return true, nil
}
