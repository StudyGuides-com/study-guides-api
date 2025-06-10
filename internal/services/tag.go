package services

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
)

type TagService struct {
	tagpb.UnimplementedTagServiceServer
}

func NewTagService() *TagService {
	return &TagService{}
}

func (s *TagService) GetTag(ctx context.Context, req *tagpb.GetTagRequest) (*sharedpb.Tag, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &sharedpb.Tag{
			Id:            req.Id,
			Name:          "Placeholder Tag",
			Type:          sharedpb.TagType_TAG_TYPE_UNSPECIFIED,
			Context:       "placeholder",
			ContentRating: sharedpb.ContentRating_CONTENT_RATING_UNSPECIFIED,
			Public:        true,
			AccessCount:   0,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*sharedpb.Tag), nil
}

func (s *TagService) ListTagsByParent(ctx context.Context, req *tagpb.ListTagsByParentRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &tagpb.ListTagsResponse{
			Tags: []*sharedpb.Tag{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}

func (s *TagService) ListTagsByType(ctx context.Context, req *tagpb.ListTagsByTypeRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &tagpb.ListTagsResponse{
			Tags: []*sharedpb.Tag{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}

func (s *TagService) ListRootTags(ctx context.Context, req *tagpb.ListRootTagsRequest) (*tagpb.ListTagsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &tagpb.ListTagsResponse{
			Tags: []*sharedpb.Tag{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}
