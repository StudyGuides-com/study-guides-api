package services

import (
	"context"
	"log"

	interactionpb "github.com/studyguides-com/study-guides-api/api/v1/interaction"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	switch req.InteractionType {
	case sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_CORRECTLY:
		return s.answerCorrectly(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_INCORRECTLY:
		return s.answerIncorrectly(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_EASY:
		return s.answerEasy(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_HARD:
		return s.answerHard(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_REVEAL:
		return s.reveal(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_VIEW_LEARN_MORE:
		return s.viewLearnMore(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_VIEW_PASSAGE:
		return s.viewPassage(ctx, req)
	case sharedpb.InteractionType_INTERACTION_TYPE_NONE:
		return s.none(ctx, req)
	}
	return nil, status.Error(codes.InvalidArgument, "invalid interaction type")
}

func (s *InteractionService) answerCorrectly(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("answerCorrectly %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	question, err := s.store.InteractionStore().AnswerCorrectly(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{
		Question: question,
	}, nil
}

func (s *InteractionService) answerIncorrectly(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("answerIncorrectly %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	question, err := s.store.InteractionStore().AnswerIncorrectly(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{
		Question: question,
	}, nil
}

func (s *InteractionService) answerEasy(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("answerEasy %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	question, err := s.store.InteractionStore().AnswerEasy(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{
		Question: question,
	}, nil
}

func (s *InteractionService) answerHard(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("answerHard %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	question, err := s.store.InteractionStore().AnswerHard(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{
		Question: question,
	}, nil
}

func (s *InteractionService) reveal(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("reveal %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	err := s.store.InteractionStore().Reveal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{}, nil
}

func (s *InteractionService) viewLearnMore(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("viewLearnMore %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	err := s.store.InteractionStore().ViewLearnMore(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{}, nil
}

func (s *InteractionService) viewPassage(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("viewPassage %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	err := s.store.InteractionStore().ViewPassage(ctx, req)
	if err != nil {
		return nil, err
	}
	return &interactionpb.InteractResponse{}, nil
}

func (s *InteractionService) none(ctx context.Context, req *interactionpb.InteractRequest) (*interactionpb.InteractResponse, error) {
	log.Printf("none %s %s %s", req.StudyMethod, req.InteractionType, req.DeckAssignment)
	return &interactionpb.InteractResponse{}, nil
}