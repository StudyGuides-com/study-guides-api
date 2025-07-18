package services

import (
	"context"
	"log"
	"time"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	userpb "github.com/studyguides-com/study-guides-api/api/v1/user"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Profile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			log.Printf("Profile request from anonymous user")
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		log.Printf("Profile request from user %s", *session.UserID)

		// TODO: Implement actual user profile lookup
		// For now, return a placeholder response
		name := "Test User"
		gamerTag := "test_gamer123"
		email := "test@example.com"
		image := "https://example.com/avatars/test.png"
		contentTagId := "content_123"

		return &userpb.ProfileResponse{
			User: &sharedpb.User{
				Id:            *session.UserID,
				Name:          &name,
				GamerTag:      &gamerTag,
				Email:         &email,
				EmailVerified: timestamppb.New(time.Now().Add(-24 * time.Hour)), // Verified 24 hours ago
				Image:         &image,
				ContentTagId:  &contentTagId,
			},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*userpb.ProfileResponse), nil
}
