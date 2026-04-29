package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Marwan051/final_project_backend/internal/service/db_tools/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type DbToolsClient struct {
	client         pb.DbToolsServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewDbToolsClient(cfg ClientConfig) (*DbToolsClient, error) {
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

	return &DbToolsClient{
		client:         pb.NewDbToolsServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *DbToolsClient) Close() error {
	return c.conn.Close()
}

func (c *DbToolsClient) NearbyTrips(ctx context.Context, req *pb.NearbyTripsRequest) (*pb.NearbyTripsResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.NearbyTrips(ctx, req)
}

func (c *DbToolsClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return true, nil
}
