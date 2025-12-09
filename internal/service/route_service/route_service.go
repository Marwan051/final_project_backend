package route_service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Marwan051/final_project_backend/internal/service/route_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// --- 1. Domain Models (Pure Go, no Proto dependency) ---

type FindRouteParams struct {
	StartLat      float64
	StartLon      float64
	EndLat        float64
	EndLon        float64
	MaxTransfers  int32
	WalkingCutoff float64
}

type JourneyCosts struct {
	Money         float64 `json:"money"`
	TransportTime float64 `json:"transport_time"`
	Walk          float64 `json:"walk"`
}

type Journey struct {
	Path  []string     `json:"path"`
	Costs JourneyCosts `json:"costs"`
}

// --- 2. The Interface (Contract) ---

// Router defines what the service can do.
// Handlers will depend on this Interface, not the concrete struct.
type Router interface {
	FindRoute(ctx context.Context, params FindRouteParams) ([]Journey, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}

// --- 3. The Concrete Implementation ---

type grpcRouteService struct {
	client pb.RoutingServiceClient
	conn   *grpc.ClientConn
}

// NewRouteService returns the Interface, not the struct
func NewRouteService(addr string) (Router, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to routing service: %w", err)
	}

	client := pb.NewRoutingServiceClient(conn)

	return &grpcRouteService{
		client: client,
		conn:   conn,
	}, nil
}

func (s *grpcRouteService) Close() error {
	return s.conn.Close()
}

func (s *grpcRouteService) FindRoute(ctx context.Context, params FindRouteParams) ([]Journey, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Map Domain -> Proto
	req := &pb.RouteRequest{
		StartLat:      params.StartLat,
		StartLon:      params.StartLon,
		EndLat:        params.EndLat,
		EndLon:        params.EndLon,
		MaxTransfers:  params.MaxTransfers,
		WalkingCutoff: params.WalkingCutoff,
	}

	resp, err := s.client.FindRoute(ctx, req)
	if err != nil {
		return nil, err
	}

	// Map Proto -> Domain
	journeys := make([]Journey, len(resp.Journeys))
	for i, j := range resp.GetJourneys() {
		var costs JourneyCosts
		if c := j.GetCosts(); c != nil {
			costs = JourneyCosts{
				Money:         c.GetMoney(),
				TransportTime: c.GetTransportTime(),
				Walk:          c.GetWalk(),
			}
		}

		journeys[i] = Journey{
			Path:  j.GetPath(),
			Costs: costs,
		}
	}

	return journeys, nil
}

func (s *grpcRouteService) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		return false, err
	}
	return true, nil
}
