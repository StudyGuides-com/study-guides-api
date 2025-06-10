package services

import (
	"context"
	"log"

	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
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
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string) (interface{}, error) {
		log.Printf("Search request from user %s: query=%s", *userID, req.Query)

		opts := &search.SearchOptions{
			UserID:      userID,
			ContextType: FromProtoContextType(req.Context),
		}
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
