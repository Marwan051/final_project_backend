package route_service

import (
	"context"
)

const (
	DefaultMaxTransfers  int32   = 3
	DefaultWalkingCutoff float64 = 500
)

type RouteRequest struct {
	StartLat      float64 `json:"start_lat"`
	StartLon      float64 `json:"start_lon"`
	EndLat        float64 `json:"end_lat"`
	EndLon        float64 `json:"end_lon"`
	MaxTransfers  int32   `json:"max_transfers"`
	WalkingCutoff float64 `json:"walking_cutoff"`
}

// ApplyDefaults sets default values for optional fields
func (r *RouteRequest) ApplyDefaults() {
	if r.MaxTransfers == 0 {
		r.MaxTransfers = DefaultMaxTransfers
	}
	if r.WalkingCutoff == 0 {
		r.WalkingCutoff = DefaultWalkingCutoff
	}
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

type RouteResponse struct {
	NumJourneys     int       `json:"num_journeys"`
	Journeys        []Journey `json:"journeys"`
	StartTripsFound int       `json:"start_trips_found"`
	EndTripsFound   int       `json:"end_trips_found"`
}

type Router interface {
	FindRoute(ctx context.Context, req RouteRequest) (RouteResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
