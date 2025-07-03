package devops

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	devopspb "github.com/studyguides-com/study-guides-api/api/v1/devops"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DigitalOceanDevopsStore struct {
	token string
	client *godo.Client
}

func (d *DigitalOceanDevopsStore) Deploy(ctx context.Context, appID string) (string, error) {
	if d.client == nil {
		d.client = godo.NewFromToken(d.token)
	}

	// Create a new deployment for the app
	deployment, _, err := d.client.Apps.CreateDeployment(ctx, appID, &godo.DeploymentCreateRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create deployment: %w", err)
	}

	return deployment.ID, nil
}

func (d *DigitalOceanDevopsStore) Rollback(ctx context.Context, appID string, deploymentID string) (string, error) {
	if d.client == nil {
		d.client = godo.NewFromToken(d.token)
	}

	// Create a rollback deployment
	deployment, _, err := d.client.Apps.CreateDeployment(ctx, appID, &godo.DeploymentCreateRequest{
		ForceBuild: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create rollback deployment: %w", err)
	}

	return deployment.ID, nil
}

func (d *DigitalOceanDevopsStore) ListDeployments(ctx context.Context, appID string) ([]devopspb.Deployment, error) {
	if d.client == nil {
		d.client = godo.NewFromToken(d.token)
	}

	// List deployments for the app
	deployments, _, err := d.client.Apps.ListDeployments(ctx, appID, &godo.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	// Convert to proto format
	result := make([]devopspb.Deployment, len(deployments))
	for i, deployment := range deployments {
		result[i] = devopspb.Deployment{
			DeploymentId: deployment.ID,
			AppId:        appID,
			Status:       convertDeploymentStatus(deployment.Phase),
			CreatedAt:    timestamppb.New(deployment.CreatedAt),
		}
	}

	return result, nil
}

func (d *DigitalOceanDevopsStore) GetDeploymentStatus(ctx context.Context, appID string, deploymentID string) (devopspb.Deployment, error) {
	if d.client == nil {
		d.client = godo.NewFromToken(d.token)
	}

	// Get the specific deployment
	deployment, _, err := d.client.Apps.GetDeployment(ctx, appID, deploymentID)
	if err != nil {
		return devopspb.Deployment{}, fmt.Errorf("failed to get deployment: %w", err)
	}

	return devopspb.Deployment{
		DeploymentId: deployment.ID,
		AppId:        appID,
		Status:       convertDeploymentStatus(deployment.Phase),
		CreatedAt:    timestamppb.New(deployment.CreatedAt),
	}, nil
}

// convertDeploymentStatus converts DigitalOcean deployment phase to our proto enum
func convertDeploymentStatus(phase godo.DeploymentPhase) devopspb.DeploymentStatus {
	switch string(phase) {
	case "PENDING":
		return devopspb.DeploymentStatus_PENDING
	case "BUILDING":
		return devopspb.DeploymentStatus_BUILDING
	case "ACTIVE":
		return devopspb.DeploymentStatus_ACTIVE
	case "ERROR":
		return devopspb.DeploymentStatus_ERROR
	case "CANCELED":
		return devopspb.DeploymentStatus_CANCELED
	default:
		return devopspb.DeploymentStatus_UNKNOWN
	}
}