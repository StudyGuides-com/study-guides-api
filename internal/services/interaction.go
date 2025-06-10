package services

import (
	"context"

	interactionpb "github.com/studyguides-com/study-guides-api/api/v1/interaction"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

type InteractionService struct {
	interactionpb.UnimplementedInteractionServiceServer
	store store.Store
}

func NewInteractionService(store store.Store) *InteractionService {
	return &InteractionService{
		store: store,
	}
}

func (s *InteractionService) Interact(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	return nil, nil
}