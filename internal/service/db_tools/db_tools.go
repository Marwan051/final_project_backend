package db_tools

import (
	"context"

	pb "github.com/Marwan051/final_project_backend/internal/service/db_tools/proto"
)

type DbTools interface {
	NearbyTrips(ctx context.Context, req *pb.NearbyTripsRequest) (*pb.NearbyTripsResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
