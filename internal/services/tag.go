package services

import (
	"context"
	"log"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	tagpb "github.com/studyguides-com/study-guides-api/api/v1/tag"
)

type TagService struct {
	tagpb.UnimplementedTagServiceServer
}

func NewTagService() *TagService {
	return &TagService{}
}

func (s *TagService) GetTag(ctx context.Context, req *tagpb.GetTagRequest) (*sharedpb.Tag, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error) {
		if userID != nil {
			log.Printf("GetTag request from user %s: id=%s", *userID, req.Id)
		} else {
			log.Printf("GetTag request from anonymous user: id=%s", req.Id)
		}

		// TODO: Implement actual tag lookup
		// For now, return a placeholder response
		return &sharedpb.Tag{
			Id:            req.Id,
			Name:          "Placeholder Tag",
			Type:          "placeholder",
			Context:       "placeholder",
			ContentRating: "G",
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
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error) {
		if userID != nil {
			log.Printf("ListTagsByParent request from user %s: parent_id=%s", *userID, req.ParentId)
		} else {
			log.Printf("ListTagsByParent request from anonymous user: parent_id=%s", req.ParentId)
		}

		// TODO: Implement actual tag listing by parent
		// For now, return empty results
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
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error) {
		if userID != nil {
			log.Printf("ListTagsByType request from user %s: type=%s", *userID, req.Type)
		} else {
			log.Printf("ListTagsByType request from anonymous user: type=%s", req.Type)
		}

		// TODO: Implement actual tag listing by type
		// For now, return empty results
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
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error) {
		if userID != nil {
			log.Printf("ListRootTags request from user %s", *userID)
		} else {
			log.Printf("ListRootTags request from anonymous user")
		}

		// TODO: Implement actual root tag listing
		// For now, return empty results
		return &tagpb.ListTagsResponse{
			Tags: []*sharedpb.Tag{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*tagpb.ListTagsResponse), nil
}
