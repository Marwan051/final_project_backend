package route_service

import (
	"context"
)

const (
	DefaultMaxTransfers  int32   = 3
	DefaultWalkingCutoff float64 = 500
	DefaultTopK          int32   = 5
)

// RouteRequest represents a request to find routes between two locations
type RouteRequest struct {
	StartLat        float64         `json:"start_lat"`
	StartLon        float64         `json:"start_lon"`
	EndLat          float64         `json:"end_lat"`
	EndLon          float64         `json:"end_lon"`
	MaxTransfers    int32           `json:"max_transfers"`
	WalkingCutoff   float64         `json:"walking_cutoff"`
	RestrictedModes []string        `json:"restricted_modes,omitempty"`
	Weights         *RoutingWeights `json:"weights,omitempty"`
	TopK            int32           `json:"top_k,omitempty"`
}

// RoutingWeights for journey ranking
type RoutingWeights struct {
	Time     float64 `json:"time"`
	Cost     float64 `json:"cost"`
	Walk     float64 `json:"walk"`
	Transfer float64 `json:"transfer"`
}

// ApplyDefaults sets default values for optional fields
func (r *RouteRequest) ApplyDefaults() {
	if r.MaxTransfers == 0 {
		r.MaxTransfers = DefaultMaxTransfers
	}
	if r.WalkingCutoff == 0 {
		r.WalkingCutoff = DefaultWalkingCutoff
	}
	if r.TopK == 0 {
		r.TopK = DefaultTopK
	}
}

// RouteResponse contains all found journeys
type RouteResponse struct {
	NumJourneys      int       `json:"num_journeys"`
	Journeys         []Journey `json:"journeys"`
	StartTripsFound  int       `json:"start_trips_found"`
	EndTripsFound    int       `json:"end_trips_found"`
	TotalRoutesFound int       `json:"total_routes_found"`
	Error            string    `json:"error,omitempty"`
}

// Journey represents a single journey option
type Journey struct {
	ID          int            `json:"id"`
	TextSummary string         `json:"text_summary"`
	Summary     JourneySummary `json:"summary"`
	Legs        []Leg          `json:"legs"`
}

// JourneySummary contains summary metrics for a journey
type JourneySummary struct {
	TotalTimeMinutes      int      `json:"total_time_minutes"`
	TotalDistanceMeters   int      `json:"total_distance_meters"`
	WalkingDistanceMeters int      `json:"walking_distance_meters"`
	Transfers             int      `json:"transfers"`
	Cost                  float64  `json:"cost"`
	Modes                 []string `json:"modes"`
}

// Leg represents a segment of the journey
type Leg struct {
	Type     string       `json:"type"` // "walk", "trip", or "transfer"
	Walk     *WalkLeg     `json:"walk,omitempty"`
	Trip     *TripLeg     `json:"trip,omitempty"`
	Transfer *TransferLeg `json:"transfer,omitempty"`
}

// WalkLeg represents a walking segment
type WalkLeg struct {
	DistanceMeters  int          `json:"distance_meters"`
	DurationMinutes int          `json:"duration_minutes"`
	Path            []Coordinate `json:"path,omitempty"`
}

// TripLeg represents a transit trip segment
type TripLeg struct {
	TripID          string       `json:"trip_id"`
	Mode            string       `json:"mode"`
	RouteShortName  string       `json:"route_short_name"`
	Headsign        string       `json:"headsign"`
	Fare            float64      `json:"fare"`
	DurationMinutes int          `json:"duration_minutes"`
	From            Stop         `json:"from"`
	To              Stop         `json:"to"`
	Path            []Coordinate `json:"path,omitempty"`
}

// TransferLeg represents a transfer between trips
type TransferLeg struct {
	FromTripID            string       `json:"from_trip_id"`
	ToTripID              string       `json:"to_trip_id"`
	FromTripName          string       `json:"from_trip_name"`
	ToTripName            string       `json:"to_trip_name"`
	WalkingDistanceMeters int          `json:"walking_distance_meters"`
	DurationMinutes       int          `json:"duration_minutes"`
	Path                  []Coordinate `json:"path,omitempty"`
}

// Stop represents a transit stop
type Stop struct {
	StopID int        `json:"stop_id"`
	Name   string     `json:"name"`
	Coord  Coordinate `json:"coord"`
}

// Coordinate represents a geographic coordinate
type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Router interface for route finding services
type Router interface {
	FindRoute(ctx context.Context, req RouteRequest) (RouteResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
