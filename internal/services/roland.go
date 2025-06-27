package services

import (
	"context"

	rolandpb "github.com/studyguides-com/study-guides-api/api/v1/roland"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

type RolandService struct {
	rolandpb.UnimplementedRolandServiceServer
	store store.Store
}

func NewRolandService(store store.Store) *RolandService {
	return &RolandService{
		store: store,
	}
}

func (s *RolandService) SaveBundle(ctx context.Context, req *rolandpb.SaveBundleRequest) (*rolandpb.SaveBundleResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		created, err := s.store.RolandStore().SaveBundle(ctx, req.Bundle, req.Force)
		if err != nil {
			return nil, err
		}
		return &rolandpb.SaveBundleResponse{
			Created: created,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.SaveBundleResponse), nil
}

func (s *RolandService) Bundles(ctx context.Context, req *rolandpb.BundlesRequest) (*rolandpb.BundlesResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		bundles, err := s.store.RolandStore().Bundles(ctx)
		if err != nil {
			return nil, err
		}
		return &rolandpb.BundlesResponse{
			Bundles: bundles,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.BundlesResponse), nil
}

func (s *RolandService) BundlesByParserType(ctx context.Context, req *rolandpb.BundlesByParserTypeRequest) (*rolandpb.BundlesByParserTypeResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Convert the proto enum to the shared enum
		parserType := sharedpb.ParserType(req.ParserType)
		bundles, err := s.store.RolandStore().BundlesByParserType(ctx, parserType)
		if err != nil {
			return nil, err
		}
		return &rolandpb.BundlesByParserTypeResponse{
			Bundles: bundles,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.BundlesByParserTypeResponse), nil
}

func (s *RolandService) UpdateGob(ctx context.Context, req *rolandpb.UpdateGobRequest) (*rolandpb.UpdateGobResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		updated, err := s.store.RolandStore().UpdateGob(ctx, req.Id, req.GobPayload, req.Force)
		if err != nil {
			return nil, err
		}
		return &rolandpb.UpdateGobResponse{
			Updated: updated,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.UpdateGobResponse), nil
}

func (s *RolandService) DeleteAllBundles(ctx context.Context, req *rolandpb.DeleteAllBundlesRequest) (*rolandpb.DeleteAllBundlesResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		err := s.store.RolandStore().DeleteAllBundles(ctx)
		if err != nil {
			return nil, err
		}
		return &rolandpb.DeleteAllBundlesResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.DeleteAllBundlesResponse), nil
}

func (s *RolandService) DeleteBundleByID(ctx context.Context, req *rolandpb.DeleteBundleByIDRequest) (*rolandpb.DeleteBundleByIDResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		err := s.store.RolandStore().DeleteBundleByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		return &rolandpb.DeleteBundleByIDResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.DeleteBundleByIDResponse), nil
}

func (s *RolandService) DeleteBundlesByShortID(ctx context.Context, req *rolandpb.DeleteBundlesByShortIDRequest) (*rolandpb.DeleteBundlesByShortIDResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		deletedCount, err := s.store.RolandStore().DeleteBundlesByShortID(ctx, req.ShortId)
		if err != nil {
			return nil, err
		}
		return &rolandpb.DeleteBundlesByShortIDResponse{
			DeletedCount: int32(deletedCount),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.DeleteBundlesByShortIDResponse), nil
}

func (s *RolandService) MarkBundleExported(ctx context.Context, req *rolandpb.MarkBundleExportedRequest) (*rolandpb.MarkBundleExportedResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Convert the proto enum to the shared enum
		exportType := sharedpb.ExportType(req.ExportType)
		updated, err := s.store.RolandStore().MarkBundleExported(ctx, req.Id, exportType)
		if err != nil {
			return nil, err
		}
		return &rolandpb.MarkBundleExportedResponse{
			Updated: updated,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*rolandpb.MarkBundleExportedResponse), nil
}
