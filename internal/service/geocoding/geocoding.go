package geocoding

import (
	"context"

	pb "github.com/Marwan051/final_project_backend/internal/service/geocoding/proto"
)

type Geocoding interface {
	Geocode(ctx context.Context, req *pb.GeocodeRequest) (*pb.GeocodeResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	Close() error
}
