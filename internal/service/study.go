package service

import (
	"context"

	studypb "github.com/studyguides-com/study-guides-api/api/study"
	healthpb "github.com/studyguides-com/study-guides-api/api/study/health"
)

type StudyService struct {
	studypb.UnimplementedStudyServiceServer
}

func NewStudyService() *StudyService {
	return &StudyService{}
}

func (s *StudyService) HealthCheck(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthStatus_SERVING,
	}, nil
}
