package services

import (
	"context"
	"log"

	healthpb "github.com/studyguides-com/study-guides-api/api/v1/health"
)

type HealthService struct {
	healthpb.UnimplementedHealthServiceServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Check(ctx context.Context, req *healthpb.CheckRequest) (*healthpb.CheckResponse, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error) {
		if userID != nil {
			log.Printf("Logged in as %s", *userID)
		}
		return &healthpb.CheckResponse{Status: healthpb.HealthStatus_SERVING}, nil
	})
	return resp.(*healthpb.CheckResponse), err
}
