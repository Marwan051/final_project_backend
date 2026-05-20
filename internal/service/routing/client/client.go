package client

import (
	"context"
	"fmt"
	"time"

	authpkg "github.com/Marwan051/final_project_backend/internal/server"
	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
	pb "github.com/Marwan051/final_project_backend/internal/service/routing/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	pbReq := &pb.JourneyRequest{
		StartLat:      req.StartLat,
		StartLon:      req.StartLon,
		EndLat:        req.EndLat,
		EndLon:        req.EndLon,
		MaxTransfers:  req.MaxTransfers,
		WalkingCutoff: req.WalkingCutoff,
		Priority:      req.Priority,
		TopK:          req.TopK,
		Filters:       mapFilters(req.Filters),
	}

	if len(req.Weights) > 0 {
		pbReq.Weights = copyFloat64Map(req.Weights)
	}

	resp, err := c.client.FindJourneys(ctx, pbReq)
	if err != nil {
		return route_service.RouteResponse{}, fmt.Errorf("grpc findjourneys failed: %w", err)
	}

	return mapProtoToDomain(resp), nil
}

func (c *RoutingClient) HealthCheck(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	_, err := c.client.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *RoutingClient) ReloadPrefixTimes(ctx context.Context) (route_service.AdminOperationResponse, error) {
	if token, ok := authpkg.GetAuthToken(ctx); ok && token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	resp, err := c.client.ReloadPrefixTimes(ctx, &pb.Empty{})
	if err != nil {
		return route_service.AdminOperationResponse{}, fmt.Errorf("grpc reloadprefixtimes failed: %w", err)
	}

	return mapAdminOperationResponse(resp), nil
}

func (c *RoutingClient) RebuildNetwork(ctx context.Context) (route_service.AdminOperationResponse, error) {
	if token, ok := authpkg.GetAuthToken(ctx); ok && token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	resp, err := c.client.RebuildNetwork(ctx, &pb.Empty{})
	if err != nil {
		return route_service.AdminOperationResponse{}, fmt.Errorf("grpc rebuildnetwork failed: %w", err)
	}

	return mapAdminOperationResponse(resp), nil
}

func mapFilterBlock(block *route_service.FilterBlock) *pb.FilterBlock {
	if block == nil {
		return nil
	}

	return &pb.FilterBlock{
		Include:      append([]string(nil), block.Include...),
		Exclude:      append([]string(nil), block.Exclude...),
		IncludeMatch: block.IncludeMatch,
	}
}

func mapFilters(filters *route_service.Filters) *pb.Filters {
	if filters == nil {
		return nil
	}

	pbFilters := &pb.Filters{}
	pbFilters.Modes = mapFilterBlock(filters.Modes)
	pbFilters.MainStreets = mapFilterBlock(filters.MainStreets)

	if pbFilters.Modes == nil && pbFilters.MainStreets == nil {
		return nil
	}

	return pbFilters
}

func copyFloat64Map(values map[string]float64) map[string]float64 {
	if len(values) == 0 {
		return nil
	}

	result := make(map[string]float64, len(values))
	for key, value := range values {
		result[key] = value
	}

	return result
}

func mapAdminOperationResponse(resp *pb.AdminOperationResponse) route_service.AdminOperationResponse {
	if resp == nil {
		return route_service.AdminOperationResponse{}
	}

	return route_service.AdminOperationResponse{
		Status:        resp.Status,
		Message:       resp.Message,
		TripsReloaded: resp.TripsReloaded,
	}
}

func mapStopInfo(info *pb.StopInfo) *route_service.StopInfo {
	if info == nil {
		return nil
	}

	return &route_service.StopInfo{
		StopID: info.StopId,
		Name:   info.Name,
		NameAr: info.NameAr,
		Coord:  append([]float64(nil), info.Coord...),
	}
}

func mapProtoToDomain(resp *pb.JourneyResponse) route_service.RouteResponse {
	if resp == nil {
		return route_service.RouteResponse{}
	}

	journeys := make([]route_service.Journey, 0, len(resp.Journeys))
	for _, pj := range resp.Journeys {
		legs := make([]route_service.Leg, 0, len(pj.Legs))
		for _, pl := range pj.Legs {
			legs = append(legs, route_service.Leg{
				Type:                  pl.Type,
				DistanceMeters:        pl.DistanceMeters,
				DurationMinutes:       pl.DurationMinutes,
				Polyline:              pl.Polyline,
				TripID:                pl.TripId,
				TripIDs:               append([]string(nil), pl.TripIds...),
				ModeEn:                pl.ModeEn,
				ModeAr:                pl.ModeAr,
				RouteShortName:        pl.RouteShortName,
				RouteShortNameAr:      pl.RouteShortNameAr,
				Headsign:              pl.Headsign,
				HeadsignAr:            pl.HeadsignAr,
				Fare:                  pl.Fare,
				FromStop:              mapStopInfo(pl.FromStop),
				ToStop:                mapStopInfo(pl.ToStop),
				FromTripID:            pl.FromTripId,
				ToTripID:              pl.ToTripId,
				FromTripName:          pl.FromTripName,
				FromTripNameAr:        pl.FromTripNameAr,
				ToTripName:            pl.ToTripName,
				ToTripNameAr:          pl.ToTripNameAr,
				EndStopID:             pl.EndStopId,
				WalkingDistanceMeters: pl.WalkingDistanceMeters,
			})
		}

		var summary *route_service.JourneySummary
		if pj.Summary != nil {
			summary = &route_service.JourneySummary{
				TotalTimeMinutes:      pj.Summary.TotalTimeMinutes,
				WalkingDistanceMeters: pj.Summary.WalkingDistanceMeters,
				TransitDistanceMeters: pj.Summary.TransitDistanceMeters,
				TotalDistanceMeters:   pj.Summary.TotalDistanceMeters,
				Transfers:             pj.Summary.Transfers,
				Cost:                  pj.Summary.Cost,
				ModesEn:               append([]string(nil), pj.Summary.ModesEn...),
				ModesAr:               append([]string(nil), pj.Summary.ModesAr...),
				MainStreetsEn:         append([]string(nil), pj.Summary.MainStreetsEn...),
				MainStreetsAr:         append([]string(nil), pj.Summary.MainStreetsAr...),
			}
		}

		journeys = append(journeys, route_service.Journey{
			ID:             pj.Id,
			TextSummary:    pj.TextSummary,
			TextSummaryEn:  pj.TextSummaryEn,
			Summary:        summary,
			Legs:           legs,
			Labels:         append([]string(nil), pj.Labels...),
			LabelsAr:       append([]string(nil), pj.LabelsAr...),
			RecommendedFor: pj.RecommendedFor,
		})
	}

	return route_service.RouteResponse{
		GeometryEncoding: resp.GeometryEncoding,
		SelectedPriority: resp.SelectedPriority,
		WeightsUsed:      copyFloat64Map(resp.WeightsUsed),
		NumJourneys:      resp.NumJourneys,
		Journeys:         journeys,
		StartTripsFound:  resp.StartTripsFound,
		EndTripsFound:    resp.EndTripsFound,
		TotalRoutesFound: resp.TotalRoutesFound,
		TotalAfterDedup:  resp.TotalAfterDedup,
		Error:            resp.Error,
	}
}
