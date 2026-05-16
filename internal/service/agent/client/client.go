package client

import (
	"context"
	"fmt"
	"time"

	agent_service "github.com/Marwan051/final_project_backend/internal/service/agent"
	pb "github.com/Marwan051/final_project_backend/internal/service/agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (c *AgentClient) Query(ctx context.Context, req agent_service.Request) (agent_service.Response, error) {
	if err := ctx.Err(); err != nil {
		return agent_service.Response{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	resp, err := c.client.Query(ctx, &pb.QueryRequest{
		Message:   req.Message,
		SessionId: req.SessionID,
	})
	if err != nil {
		return agent_service.Response{}, fmt.Errorf("grpc query failed: %w", err)
	}

	return agent_service.Response{
		Message:   resp.Message,
		SessionID: resp.SessionId,
	}, nil
}

func (c *AgentClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		return false, err
	}
	return true, nil
}
