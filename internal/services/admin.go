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

func (s *AdminService) KillUser(ctx context.Context, req *adminpb.KillUserAdminRequest) (*adminpb.KillUserAdminResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			log.Printf("KillUser request from anonymous user")
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		log.Printf("KillUser request from user %s for email %s", *session.UserID, req.Email)

		// Call the store method to kill the user
		ok, err := s.store.UserStore().KillUser(ctx, req.Email)
		if err != nil {
			log.Printf("Error killing user %s: %v", req.Email, err)
			return nil, status.Error(codes.Internal, "failed to kill user")
		}

		return &adminpb.KillUserAdminResponse{
			Ok: ok,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*adminpb.KillUserAdminResponse), nil
}
