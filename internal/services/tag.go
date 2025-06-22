package services

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TagService struct {
	tagpb.UnimplementedTagServiceServer
	store store.Store
}

func NewTagService(store store.Store) *TagService {
	return &TagService{
		store: store,
	}
}

func (s *TagService) GetTag(ctx context.Context, req *tagpb.GetTagRequest) (*sharedpb.Tag, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return s.store.TagStore().GetTagByID(ctx, req.Id)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*sharedpb.Tag), nil
}

func (s *TagService) ListTagsByParent(ctx context.Context, req *tagpb.ListTagsByParentRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		tags, err := s.store.TagStore().ListTagsByParent(ctx, req.ParentId)
		if err != nil {
			return nil, err
		}
		return &tagpb.ListTagsResponse{
			Tags: tags,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}

func (s *TagService) ListTagsByType(ctx context.Context, req *tagpb.ListTagsByTypeRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		tags, err := s.store.TagStore().ListTagsByType(ctx, sharedpb.TagType(sharedpb.TagType_value[req.Type]))
		if err != nil {
			return nil, err
		}
		return &tagpb.ListTagsResponse{
			Tags: tags,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}

func (s *TagService) ListRootTags(ctx context.Context, req *tagpb.ListRootTagsRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		tags, err := s.store.TagStore().ListRootTags(ctx)
		if err != nil {
			return nil, err
		}
		return &tagpb.ListTagsResponse{
			Tags: tags,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}

func (s *TagService) Report(ctx context.Context, req *tagpb.ReportTagRequest) (*tagpb.ReportTagResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "user must be authenticated to report tags")
		}
		err := s.store.TagStore().Report(ctx, req.TagId, *session.UserID, req.ReportType, req.Reason)
		if err != nil {
			return nil, err
		}
		return &tagpb.ReportTagResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ReportTagResponse), nil
}

func (s *TagService) Favorite(ctx context.Context, req *tagpb.FavoriteTagRequest) (*tagpb.FavoriteTagResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "user must be authenticated to favorite tags")
		}
		err := s.store.TagStore().Favorite(ctx, req.TagId, *session.UserID)
		if err != nil {
			return nil, err
		}
		return &tagpb.FavoriteTagResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.FavoriteTagResponse), nil
}

func (s *TagService) Unfavorite(ctx context.Context, req *tagpb.UnfavoriteTagRequest) (*tagpb.UnfavoriteTagResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		if session.UserID == nil {
			return nil, status.Error(codes.Unauthenticated, "user must be authenticated to unfavorite tags")
		}
		err := s.store.TagStore().Unfavorite(ctx, req.TagId, *session.UserID)
		if err != nil {
			return nil, err
		}
		return &tagpb.UnfavoriteTagResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.UnfavoriteTagResponse), nil
}
