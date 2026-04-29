package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
	pb "github.com/Marwan051/final_project_backend/internal/service/routing/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address        string
	RequestTimeout time.Duration
	DialOptions    []grpc.DialOption
}

type RoutingClient struct {
	client         pb.RoutingServiceClient
	conn           *grpc.ClientConn
	requestTimeout time.Duration
}

func NewRoutingClient(cfg ClientConfig) (*RoutingClient, error) {
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

	return &RoutingClient{
		client:         pb.NewRoutingServiceClient(conn),
		conn:           conn,
		requestTimeout: timeout,
	}, nil
}

func (c *RoutingClient) Close() error {
	return c.conn.Close()
}

func (c *RoutingClient) FindRoute(ctx context.Context, req route_service.RouteRequest) (route_service.RouteResponse, error) {
	if err := ctx.Err(); err != nil {
		return route_service.RouteResponse{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req.ApplyDefaults()

	var excludeModes []string
	if req.RestrictedModes != nil {
		excludeModes = req.RestrictedModes
	}

	filters := &pb.Filters{}
	if len(excludeModes) > 0 {
		filters.Modes = &pb.FilterBlock{
			Exclude: excludeModes,
		}
	}

	pbReq := &pb.JourneyRequest{
		StartLat:      req.StartLat,
		StartLon:      req.StartLon,
		EndLat:        req.EndLat,
		EndLon:        req.EndLon,
		MaxTransfers:  req.MaxTransfers,
		WalkingCutoff: int32(req.WalkingCutoff),
		TopK:          req.TopK,
		Filters:       filters,
	}

	if req.Weights != nil {
		pbReq.Weights = map[string]float64{
			"time":     req.Weights.Time,
			"cost":     req.Weights.Cost,
			"walk":     req.Weights.Walk,
			"transfer": req.Weights.Transfer,
		}
	}

	resp, err := c.client.FindJourneys(ctx, pbReq)
	if err != nil {
		return route_service.RouteResponse{}, fmt.Errorf("grpc findjourneys failed: %w", err)
	}

	return mapProtoToDomain(resp), nil
}

func (c *RoutingClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func decodePolyline(encoded string) []route_service.Coordinate {
	// A placeholder for polyline decoding if the domain model requires it
	return []route_service.Coordinate{}
}

func parseStop(info *pb.StopInfo) route_service.Stop {
	if info == nil {
		return route_service.Stop{}
	}
	id, _ := strconv.Atoi(info.StopId)
	coord := route_service.Coordinate{}
	if len(info.Coord) >= 2 {
		coord.Lon = info.Coord[0]
		coord.Lat = info.Coord[1]
	}
	return route_service.Stop{
		StopID: id,
		Name:   info.Name,
		Coord:  coord,
	}
}

func mapProtoToDomain(resp *pb.JourneyResponse) route_service.RouteResponse {
	journeys := make([]route_service.Journey, 0, len(resp.Journeys))
	for _, pj := range resp.Journeys {
		legs := make([]route_service.Leg, 0, len(pj.Legs))
		for _, pl := range pj.Legs {
			leg := route_service.Leg{
				Type: pl.Type,
			}

			switch pl.Type {
			case "walk":
				leg.Walk = &route_service.WalkLeg{
					DistanceMeters:  int(pl.DistanceMeters),
					DurationMinutes: int(pl.DurationMinutes),
					Path:            decodePolyline(pl.Polyline),
				}
			case "trip":
				leg.Trip = &route_service.TripLeg{
					TripID:          pl.TripId,
					Mode:            pl.ModeEn,
					RouteShortName:  pl.RouteShortName,
					Headsign:        pl.Headsign,
					Fare:            pl.Fare,
					DurationMinutes: int(pl.DurationMinutes),
					From:            parseStop(pl.FromStop),
					To:              parseStop(pl.ToStop),
					Path:            decodePolyline(pl.Polyline),
				}
			case "transfer":
				leg.Transfer = &route_service.TransferLeg{
					FromTripID:            pl.FromTripId,
					ToTripID:              pl.ToTripId,
					FromTripName:          pl.FromTripName,
					ToTripName:            pl.ToTripName,
					WalkingDistanceMeters: int(pl.WalkingDistanceMeters),
					DurationMinutes:       int(pl.DurationMinutes),
					Path:                  decodePolyline(pl.Polyline),
				}
			}
			legs = append(legs, leg)
		}

		summary := route_service.JourneySummary{
			TotalTimeMinutes:      int(pj.Summary.TotalTimeMinutes),
			TotalDistanceMeters:   int(pj.Summary.TotalDistanceMeters),
			WalkingDistanceMeters: int(pj.Summary.WalkingDistanceMeters),
			Transfers:             int(pj.Summary.Transfers),
			Cost:                  pj.Summary.Cost,
			Modes:                 pj.Summary.ModesEn,
		}

		journeys = append(journeys, route_service.Journey{
			ID:          int(pj.Id),
			TextSummary: pj.TextSummaryEn,
			Summary:     summary,
			Legs:        legs,
		})
	}

	return route_service.RouteResponse{
		NumJourneys:      int(resp.NumJourneys),
		Journeys:         journeys,
		StartTripsFound:  int(resp.StartTripsFound),
		EndTripsFound:    int(resp.EndTripsFound),
		TotalRoutesFound: int(resp.TotalRoutesFound),
		Error:            resp.Error,
	}
}
