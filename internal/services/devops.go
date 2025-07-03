package services

import (
	"context"

	devopspb "github.com/studyguides-com/study-guides-api/api/v1/devops"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevopsService struct {
	devopspb.UnimplementedDevopsServiceServer
	store store.Store
}

func NewDevopsService(store store.Store) *DevopsService {
	return &DevopsService{
		store: store,
	}
}

func (s *DevopsService) Deploy(ctx context.Context, req *devopspb.DeployRequest) (*devopspb.DeployResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		deploymentID, err := s.store.DevopsStore().Deploy(ctx, req.AppId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &devopspb.DeployResponse{DeploymentId: deploymentID}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*devopspb.DeployResponse), nil
}

func (s *DevopsService) Rollback(ctx context.Context, req *devopspb.RollbackRequest) (*devopspb.RollbackResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		deploymentID, err := s.store.DevopsStore().Rollback(ctx, req.AppId, req.DeploymentId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &devopspb.RollbackResponse{DeploymentId: deploymentID}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*devopspb.RollbackResponse), nil
}

func (s *DevopsService) ListDeployments(ctx context.Context, req *devopspb.ListDeploymentsRequest) (*devopspb.ListDeploymentsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		deployments, err := s.store.DevopsStore().ListDeployments(ctx, req.AppId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Convert slice of Deployment to slice of *Deployment
		deploymentPtrs := make([]*devopspb.Deployment, len(deployments))
		for i := range deployments {
			deploymentPtrs[i] = &deployments[i]
		}

		return &devopspb.ListDeploymentsResponse{Deployments: deploymentPtrs}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*devopspb.ListDeploymentsResponse), nil
}

func (s *DevopsService) GetDeploymentStatus(ctx context.Context, req *devopspb.GetDeploymentStatusRequest) (*devopspb.GetDeploymentStatusResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		deployment, err := s.store.DevopsStore().GetDeploymentStatus(ctx, req.AppId, req.DeploymentId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &devopspb.GetDeploymentStatusResponse{Deployment: &deployment}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*devopspb.GetDeploymentStatusResponse), nil
}