package services

import (
	"context"
	"log"

	adminpb "github.com/studyguides-com/study-guides-api/api/v1/admin"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService struct {
	adminpb.UnimplementedAdminServiceServer
	store store.Store
}

func NewAdminService(store store.Store) *AdminService {
	return &AdminService{
		store: store,
	}
}

func (s *AdminService) NewTag(ctx context.Context, req *adminpb.NewTagAdminRequest) (*adminpb.NewTagAdminResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			log.Printf("NewTag request from anonymous user")
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		log.Printf("NewTag request from user %s", *session.UserID)

		return &adminpb.NewTagAdminResponse{
			Tag: &sharedpb.Tag{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*adminpb.NewTagAdminResponse), nil
}
