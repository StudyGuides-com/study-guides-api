package devops

import (
	"context"
	"os"

	devopspb "github.com/studyguides-com/study-guides-api/api/v1/devops"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevopsStore interface {
	Deploy(ctx context.Context, appID string) (string, error)
	Rollback(ctx context.Context, appID string, deploymentID string) (string, error)
	ListDeployments(ctx context.Context, appID string) ([]devopspb.Deployment, error)
	GetDeploymentStatus(ctx context.Context, appID string, deploymentID string) (devopspb.Deployment, error)
}

func NewDevopsStore(ctx context.Context) (*DigitalOceanDevopsStore, error) {
	token := os.Getenv("DIGITAL_OCEAN_TOKEN")
	if token == "" {
		return nil, status.Error(codes.Internal, "DIGITAL_OCEAN_TOKEN is not set")
	}
	return &DigitalOceanDevopsStore{token: token}, nil
}
