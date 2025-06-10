package services

import (
	"context"

	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
)

type HealthService struct {
	healthpb.UnimplementedHealthServiceServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Check(ctx context.Context, req *healthpb.CheckRequest) (*healthpb.CheckResponse, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context) (interface{}, error) {
		return &healthpb.CheckResponse{Status: healthpb.HealthStatus_SERVING}, nil
	})
	return resp.(*healthpb.CheckResponse), err
}
