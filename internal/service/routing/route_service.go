package route_service

import (
	"context"
)

const (
	DefaultMaxTransfers  int32 = 3
	DefaultWalkingCutoff int32 = 500
	DefaultTopK          int32 = 5
)

// FilterBlock controls how a category is included or excluded.
type FilterBlock struct {
	Include      []string `json:"include,omitempty"`
	Exclude      []string `json:"exclude,omitempty"`
	IncludeMatch string   `json:"include_match,omitempty"`
}

// Filters groups the request filters.
type Filters struct {
	Modes       *FilterBlock `json:"modes,omitempty"`
	MainStreets *FilterBlock `json:"main_streets,omitempty"`
}

// RouteRequest represents a request to find routes between two locations.
type RouteRequest struct {
	StartLat      float64            `json:"start_lat"`
	StartLon      float64            `json:"start_lon"`
	EndLat        float64            `json:"end_lat"`
	EndLon        float64            `json:"end_lon"`
	MaxTransfers  int32              `json:"max_transfers"`
	WalkingCutoff int32              `json:"walking_cutoff"`
	Priority      string             `json:"priority,omitempty"`
	TopK          int32              `json:"top_k,omitempty"`
	Filters       *Filters           `json:"filters,omitempty"`
	Weights       map[string]float64 `json:"weights,omitempty"`
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

// StopInfo describes a stop in the routing response.
type StopInfo struct {
	StopID string    `json:"stop_id,omitempty"`
	Name   string    `json:"name,omitempty"`
	NameAr string    `json:"name_ar,omitempty"`
	Coord  []float64 `json:"coord,omitempty"`
}

// Leg represents one leg in a journey.
type Leg struct {
	Type                  string    `json:"type,omitempty"`
	DistanceMeters        int32     `json:"distance_meters,omitempty"`
	DurationMinutes       int32     `json:"duration_minutes,omitempty"`
	Polyline              string    `json:"polyline,omitempty"`
	TripID                string    `json:"trip_id,omitempty"`
	TripIDs               []string  `json:"trip_ids,omitempty"`
	ModeEn                string    `json:"mode_en,omitempty"`
	ModeAr                string    `json:"mode_ar,omitempty"`
	RouteShortName        string    `json:"route_short_name,omitempty"`
	RouteShortNameAr      string    `json:"route_short_name_ar,omitempty"`
	Headsign              string    `json:"headsign,omitempty"`
	HeadsignAr            string    `json:"headsign_ar,omitempty"`
	Fare                  float64   `json:"fare,omitempty"`
	FromStop              *StopInfo `json:"from_stop,omitempty"`
	ToStop                *StopInfo `json:"to_stop,omitempty"`
	FromTripID            string    `json:"from_trip_id,omitempty"`
	ToTripID              string    `json:"to_trip_id,omitempty"`
	FromTripName          string    `json:"from_trip_name,omitempty"`
	FromTripNameAr        string    `json:"from_trip_name_ar,omitempty"`
	ToTripName            string    `json:"to_trip_name,omitempty"`
	ToTripNameAr          string    `json:"to_trip_name_ar,omitempty"`
	EndStopID             string    `json:"end_stop_id,omitempty"`
	WalkingDistanceMeters int32     `json:"walking_distance_meters,omitempty"`
}

// JourneySummary contains summary metrics for a journey.
type JourneySummary struct {
	TotalTimeMinutes      int32    `json:"total_time_minutes,omitempty"`
	WalkingDistanceMeters int32    `json:"walking_distance_meters,omitempty"`
	TransitDistanceMeters int32    `json:"transit_distance_meters,omitempty"`
	TotalDistanceMeters   int32    `json:"total_distance_meters,omitempty"`
	Transfers             int32    `json:"transfers,omitempty"`
	Cost                  float64  `json:"cost,omitempty"`
	ModesEn               []string `json:"modes_en,omitempty"`
	ModesAr               []string `json:"modes_ar,omitempty"`
	MainStreetsEn         []string `json:"main_streets_en,omitempty"`
	MainStreetsAr         []string `json:"main_streets_ar,omitempty"`
}

// Journey represents a single journey option.
type Journey struct {
	ID             int32           `json:"id,omitempty"`
	TextSummary    string          `json:"text_summary,omitempty"`
	TextSummaryEn  string          `json:"text_summary_en,omitempty"`
	Summary        *JourneySummary `json:"summary,omitempty"`
	Legs           []Leg           `json:"legs,omitempty"`
	Labels         []string        `json:"labels,omitempty"`
	LabelsAr       []string        `json:"labels_ar,omitempty"`
	RecommendedFor string          `json:"recommended_for,omitempty"`
}

// RouteResponse contains all found journeys.
type RouteResponse struct {
	GeometryEncoding string             `json:"geometry_encoding,omitempty"`
	SelectedPriority string             `json:"selected_priority,omitempty"`
	WeightsUsed      map[string]float64 `json:"weights_used,omitempty"`
	NumJourneys      int32              `json:"num_journeys,omitempty"`
	Journeys         []Journey          `json:"journeys,omitempty"`
	StartTripsFound  int32              `json:"start_trips_found,omitempty"`
	EndTripsFound    int32              `json:"end_trips_found,omitempty"`
	TotalRoutesFound int32              `json:"total_routes_found,omitempty"`
	TotalAfterDedup  int32              `json:"total_after_dedup,omitempty"`
	Error            string             `json:"error,omitempty"`
}

// AdminOperationResponse contains the result of an admin-only routing task.
type AdminOperationResponse struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	TripsReloaded int32  `json:"trips_reloaded,omitempty"`
}

// Router interface for route finding services
type Router interface {
	FindRoute(ctx context.Context, req RouteRequest) (RouteResponse, error)
	ReloadPrefixTimes(ctx context.Context) (AdminOperationResponse, error)
	RebuildNetwork(ctx context.Context) (AdminOperationResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
