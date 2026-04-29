package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Marwan051/final_project_backend/internal/service/traffic_updater/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type TrafficUpdaterClient struct {
	client         pb.TrafficUpdateServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewTrafficUpdaterClient(cfg ClientConfig) (*TrafficUpdaterClient, error) {
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

	return &TrafficUpdaterClient{
		client:         pb.NewTrafficUpdateServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *TrafficUpdaterClient) Close() error {
	return c.conn.Close()
}

func (c *TrafficUpdaterClient) TriggerUpdate(ctx context.Context, req *pb.TriggerRequest) (*pb.UpdateResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.TriggerUpdate(ctx, req)
}

func (c *TrafficUpdaterClient) GetStatus(ctx context.Context) (*pb.StatusResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.GetStatus(ctx, &pb.Empty{})
}

func (c *TrafficUpdaterClient) UpdateTrip(ctx context.Context, req *pb.UpdateTripRequest) (*pb.UpdateResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.UpdateTrip(ctx, req)
}

func (c *TrafficUpdaterClient) StreetTraffic(ctx context.Context, req *pb.StreetTrafficRequest) (*pb.StreetTrafficResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.StreetTraffic(ctx, req)
}

func (c *TrafficUpdaterClient) ListStreets(ctx context.Context) (*pb.StreetListResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.ListStreets(ctx, &pb.Empty{})
}
