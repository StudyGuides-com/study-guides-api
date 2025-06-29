package services

import (
	"context"
	"log"

	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"github.com/studyguides-com/study-guides-api/internal/store/search"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchService struct {
	searchpb.UnimplementedSearchServiceServer
	store store.Store
}

func NewSearchService(s store.Store) *SearchService {
	return &SearchService{
		store: s,
	}
}

func (s *SearchService) SearchTags(ctx context.Context, req *searchpb.SearchTagsRequest) (*searchpb.SearchTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		log.Printf("Search request from user %s: query=%s", *session.UserID, req.Query)

		opts := search.NewSearchOptionsFromRequest(ctx, req)
		results, err := s.store.SearchStore().SearchTags(ctx, req.Query, opts)
		if err != nil {
			return nil, err
		}
		return &searchpb.SearchTagsResponse{
			Results: results,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, status.Error(codes.Internal, "search service returned nil response")
	}
	return resp.(*searchpb.SearchTagsResponse), nil
}

func (s *SearchService) SearchUsers(ctx context.Context, req *searchpb.SearchUsersRequest) (*searchpb.SearchUsersResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if !session.IsAuth {
			return nil, status.Error(codes.PermissionDenied, "you must logged in to search users")
		}
		log.Printf("Search request from user %s: query=%s", *session.UserID, req.Query)

		opts := search.NewSearchOptionsFrom(ctx)
		results, err := s.store.SearchStore().SearchUsers(ctx, req.Query, opts)
		if err != nil {
			return nil, err
		}
		return &searchpb.SearchUsersResponse{
			Results: results,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, status.Error(codes.Internal, "search service returned nil response")
	}
	return resp.(*searchpb.SearchUsersResponse), nil
}

func (s *SearchService) ListIndexes(ctx context.Context, req *searchpb.ListIndexesRequest) (*searchpb.ListIndexesResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &searchpb.ListIndexesResponse{
			Indexes: s.store.SearchStore().ListIndexes(ctx).Indexes,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, status.Error(codes.Internal, "search service returned nil response")
	}
	return resp.(*searchpb.ListIndexesResponse), nil
}
