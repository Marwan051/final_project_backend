package traffic_updater

import (
	"context"

	pb "github.com/Marwan051/final_project_backend/internal/service/traffic_updater/proto"
)

type TrafficUpdater interface {
	TriggerUpdate(ctx context.Context, req *pb.TriggerRequest) (*pb.UpdateResponse, error)
	GetStatus(ctx context.Context) (*pb.StatusResponse, error)
	UpdateTrip(ctx context.Context, req *pb.UpdateTripRequest) (*pb.UpdateResponse, error)
	StreetTraffic(ctx context.Context, req *pb.StreetTrafficRequest) (*pb.StreetTrafficResponse, error)
	ListStreets(ctx context.Context) (*pb.StreetListResponse, error)
	Close() error
}
