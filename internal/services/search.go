package services

import (
	"context"
	"log"

	"github.com/studyguides-com/study-guides-api/internal/store"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
)

type SearchService struct {
	searchpb.UnimplementedSearchServiceServer
	store store.Store
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) Search(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchResponse, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string) (interface{}, error) {
		if userID != nil {
			log.Printf("Search request from user %s: query=%s, context=%s", *userID, req.Query, req.Context)
		} else {
			log.Printf("Search request from anonymous user: query=%s, context=%s", req.Query, req.Context)
		}

		results, err := s.store.SearchStore().SearchTags(ctx, FromProtoContextType(req.Context), req.Query)
		if err != nil {
			return nil, err
		}
		return &searchpb.SearchResponse{
			Results: results,
		}, nil
	})
	return resp.(*searchpb.SearchResponse), err
}

