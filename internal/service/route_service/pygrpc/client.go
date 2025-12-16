package pygrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Marwan051/final_project_backend/internal/service/route_service"
	pb "github.com/Marwan051/final_project_backend/internal/service/route_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type Client struct {
	client         pb.RoutingServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewClient(cfg ClientConfig) (route_service.Router, error) {
	// Base options - simple for demo
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	opts = append(opts, cfg.DialOptions...)

	// Create Connection
	conn, err := grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client for %s: %w", cfg.Address, err)
	}

	// Set default timeout if not provided
	timeout := cfg.RequestTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		client:         pb.NewRoutingServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) FindRoute(ctx context.Context, req route_service.RouteRequest) (route_service.RouteResponse, error) {
	// Optimization: Check if context is already done before starting
	if err := ctx.Err(); err != nil {
		return route_service.RouteResponse{}, err
	}

	// Use configured timeout only if the parent context doesn't have a tighter deadline
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	// Apply defaults
	req.ApplyDefaults()

	// Map domain model -> protobuf
	pbReq := &pb.RouteRequest{
		StartLat:      req.StartLat,
		StartLon:      req.StartLon,
		EndLat:        req.EndLat,
		EndLon:        req.EndLon,
		MaxTransfers:  req.MaxTransfers,
		WalkingCutoff: req.WalkingCutoff,
	}

	resp, err := c.client.FindRoute(ctx, pbReq)
	if err != nil {
		return route_service.RouteResponse{}, fmt.Errorf("grpc findroute failed: %w", err)
	}

	return mapProtoToDomain(resp), nil
}

func (c *Client) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		return false, err
	}
	return true, nil
}

// mapProtoToDomain separates the mapping logic to keep the main method clean.
// This allows the compiler to inline this function potentially.
func mapProtoToDomain(resp *pb.RouteResponse) route_service.RouteResponse {
	journeys := make([]route_service.Journey, len(resp.GetJourneys()))

	for i, j := range resp.GetJourneys() {
		var costs route_service.JourneyCosts
		if c := j.GetCosts(); c != nil {
			costs = route_service.JourneyCosts{
				Money:         c.GetMoney(),
				TransportTime: c.GetTransportTime(),
				Walk:          c.GetWalk(),
			}
		}

		journeys[i] = route_service.Journey{
			Path:  j.GetPath(),
			Costs: costs,
		}
	}

	return route_service.RouteResponse{
		NumJourneys:     int(resp.GetNumJourneys()),
		Journeys:        journeys,
		StartTripsFound: int(resp.GetStartTripsFound()),
		EndTripsFound:   int(resp.GetEndTripsFound()),
	}
}
