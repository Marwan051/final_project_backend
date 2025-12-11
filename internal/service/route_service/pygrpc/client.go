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

// Client is the gRPC implementation of route_service.Router
type Client struct {
	client pb.RoutingServiceClient
	conn   *grpc.ClientConn
}

// NewClient creates a new gRPC routing client
func NewClient(addr string) (route_service.Router, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Wait for connection to be established
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to routing service at %s: %w", addr, err)
	}

	return &Client{
		client: pb.NewRoutingServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) FindRoute(ctx context.Context, req route_service.RouteRequest) (route_service.RouteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

	// Map protobuf -> domain model
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
	}, nil
}

func (c *Client) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.HealthRequest{})
	return err == nil, err
}
