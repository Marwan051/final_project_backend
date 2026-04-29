package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Marwan051/final_project_backend/internal/service/geocoding/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type GeocodingClient struct {
	client         pb.GeocodingServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewGeocodingClient(cfg ClientConfig) (*GeocodingClient, error) {
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

	return &GeocodingClient{
		client:         pb.NewGeocodingServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *GeocodingClient) Close() error {
	return c.conn.Close()
}

func (c *GeocodingClient) Geocode(ctx context.Context, req *pb.GeocodeRequest) (*pb.GeocodeResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	return c.client.Geocode(ctx, req)
}

func (c *GeocodingClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return true, nil
}
