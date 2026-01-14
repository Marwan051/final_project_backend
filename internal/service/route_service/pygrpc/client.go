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
	// Base options
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

	if err := ctx.Err(); err != nil {
		return route_service.RouteResponse{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req.ApplyDefaults()

	pbReq := &pb.RouteRequest{
		StartLat:        req.StartLat,
		StartLon:        req.StartLon,
		EndLat:          req.EndLat,
		EndLon:          req.EndLon,
		MaxTransfers:    req.MaxTransfers,
		WalkingCutoff:   req.WalkingCutoff,
		RestrictedModes: req.RestrictedModes,
		TopK:            req.TopK,
	}

	// Map weights if provided
	if req.Weights != nil {
		pbReq.Weights = &pb.RoutingWeights{
			Time:     req.Weights.Time,
			Cost:     req.Weights.Cost,
			Walk:     req.Weights.Walk,
			Transfer: req.Weights.Transfer,
		}
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

func mapProtoToDomain(resp *pb.RouteResponse) route_service.RouteResponse {
	journeys := make([]route_service.Journey, len(resp.GetJourneys()))

	for i, j := range resp.GetJourneys() {
		journeys[i] = route_service.Journey{
			ID:          int(j.GetId()),
			TextSummary: j.GetTextSummary(),
			Summary:     mapSummary(j.GetSummary()),
			Legs:        mapLegs(j.GetLegs()),
		}
	}

	return route_service.RouteResponse{
		NumJourneys:      int(resp.GetNumJourneys()),
		Journeys:         journeys,
		StartTripsFound:  int(resp.GetStartTripsFound()),
		EndTripsFound:    int(resp.GetEndTripsFound()),
		TotalRoutesFound: int(resp.GetTotalRoutesFound()),
		Error:            resp.GetError(),
	}
}

func mapSummary(s *pb.JourneySummary) route_service.JourneySummary {
	if s == nil {
		return route_service.JourneySummary{}
	}
	return route_service.JourneySummary{
		TotalTimeMinutes:      int(s.GetTotalTimeMinutes()),
		TotalDistanceMeters:   int(s.GetTotalDistanceMeters()),
		WalkingDistanceMeters: int(s.GetWalkingDistanceMeters()),
		Transfers:             int(s.GetTransfers()),
		Cost:                  s.GetCost(),
		Modes:                 s.GetModes(),
	}
}

func mapLegs(legs []*pb.Leg) []route_service.Leg {
	result := make([]route_service.Leg, len(legs))
	for i, leg := range legs {
		result[i] = mapLeg(leg)
	}
	return result
}

func mapLeg(leg *pb.Leg) route_service.Leg {
	if leg == nil {
		return route_service.Leg{}
	}

	switch v := leg.GetLegType().(type) {
	case *pb.Leg_Walk:
		return route_service.Leg{
			Type: "walk",
			Walk: mapWalkLeg(v.Walk),
		}
	case *pb.Leg_Trip:
		return route_service.Leg{
			Type: "trip",
			Trip: mapTripLeg(v.Trip),
		}
	case *pb.Leg_Transfer:
		return route_service.Leg{
			Type:     "transfer",
			Transfer: mapTransferLeg(v.Transfer),
		}
	default:
		return route_service.Leg{}
	}
}

func mapWalkLeg(w *pb.WalkLeg) *route_service.WalkLeg {
	if w == nil {
		return nil
	}
	return &route_service.WalkLeg{
		DistanceMeters:  int(w.GetDistanceMeters()),
		DurationMinutes: int(w.GetDurationMinutes()),
		Path:            mapCoordinates(w.GetPath()),
	}
}

func mapTripLeg(t *pb.TripLeg) *route_service.TripLeg {
	if t == nil {
		return nil
	}
	return &route_service.TripLeg{
		TripID:          t.GetTripId(),
		Mode:            t.GetMode(),
		RouteShortName:  t.GetRouteShortName(),
		Headsign:        t.GetHeadsign(),
		Fare:            t.GetFare(),
		DurationMinutes: int(t.GetDurationMinutes()),
		From:            mapStop(t.GetFrom()),
		To:              mapStop(t.GetTo()),
		Path:            mapCoordinates(t.GetPath()),
	}
}

func mapTransferLeg(t *pb.TransferLeg) *route_service.TransferLeg {
	if t == nil {
		return nil
	}
	return &route_service.TransferLeg{
		FromTripID:            t.GetFromTripId(),
		ToTripID:              t.GetToTripId(),
		FromTripName:          t.GetFromTripName(),
		ToTripName:            t.GetToTripName(),
		WalkingDistanceMeters: int(t.GetWalkingDistanceMeters()),
		DurationMinutes:       int(t.GetDurationMinutes()),
		Path:                  mapCoordinates(t.GetPath()),
	}
}

func mapStop(s *pb.Stop) route_service.Stop {
	if s == nil {
		return route_service.Stop{}
	}
	return route_service.Stop{
		StopID: int(s.GetStopId()),
		Name:   s.GetName(),
		Coord:  mapCoordinate(s.GetCoord()),
	}
}

func mapCoordinates(coords []*pb.Coordinate) []route_service.Coordinate {
	result := make([]route_service.Coordinate, len(coords))
	for i, c := range coords {
		result[i] = mapCoordinate(c)
	}
	return result
}

func mapCoordinate(c *pb.Coordinate) route_service.Coordinate {
	if c == nil {
		return route_service.Coordinate{}
	}
	return route_service.Coordinate{
		Lon: c.GetLon(),
		Lat: c.GetLat(),
	}
}
