package store

import (
	"context"
	"os"

	"github.com/studyguides-com/study-guides-api/internal/store/devops"
	"github.com/studyguides-com/study-guides-api/internal/store/interaction"
	"github.com/studyguides-com/study-guides-api/internal/store/kpi"
	"github.com/studyguides-com/study-guides-api/internal/store/question"
	"github.com/studyguides-com/study-guides-api/internal/store/roland"
	"github.com/studyguides-com/study-guides-api/internal/store/search"
	"github.com/studyguides-com/study-guides-api/internal/store/tag"
	"github.com/studyguides-com/study-guides-api/internal/store/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Store interface {
	SearchStore() search.SearchStore
	TagStore() tag.TagStore
	UserStore() user.UserStore
	QuestionStore() question.QuestionStore
	InteractionStore() interaction.InteractionStore
	RolandStore() roland.RolandStore
	DevopsStore() devops.DevopsStore
	KPIStore() kpi.KPIStore
}

type store struct {
	searchStore      search.SearchStore
	tagStore         tag.TagStore
	userStore        user.UserStore
	questionStore    question.QuestionStore
	interactionStore interaction.InteractionStore
	rolandStore      roland.RolandStore
	devopsStore      devops.DevopsStore
	kpiStore         kpi.KPIStore
}

func (s *store) SearchStore() search.SearchStore {
	return s.searchStore
}

func (s *store) TagStore() tag.TagStore {
	return s.tagStore
}

func (s *store) UserStore() user.UserStore {
	return s.userStore
}

func (s *store) QuestionStore() question.QuestionStore {
	return s.questionStore
}

func (s *store) InteractionStore() interaction.InteractionStore {
	return s.interactionStore
}

func (s *store) RolandStore() roland.RolandStore {
	return s.rolandStore
}

func (s *store) DevopsStore() devops.DevopsStore {
	return s.devopsStore
}

func (s *store) KPIStore() kpi.KPIStore {
	return s.kpiStore
}

func NewStore() (Store, error) {
	ctx := context.Background()
	algoliaAppID := os.Getenv("ALGOLIA_APP_ID")
	algoliaAdminAPIKey := os.Getenv("ALGOLIA_ADMIN_API_KEY")
	dbURL := os.Getenv("DATABASE_URL")
	rolandDBURL := os.Getenv("ROLAND_DATABASE_URL")

	if algoliaAppID == "" || algoliaAdminAPIKey == "" || dbURL == "" {
		return nil, status.Error(codes.FailedPrecondition, "missing required environment variables")
	}

	searchStore, err := search.NewAlgoliaSearchStore(ctx, algoliaAppID, algoliaAdminAPIKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userStore, err := user.NewSqlUserStore(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tagStore, err := tag.NewSqlTagStore(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	questionStore, err := question.NewSqlQuestionStore(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	interactionStore, err := interaction.NewSqlInteractionStore(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rolandStore, err := roland.NewSqlRolandStore(ctx, rolandDBURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	devopsStore, err := devops.NewDevopsStore(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	kpiStore, err := kpi.NewSqlKPIStore(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &store{
		searchStore:      searchStore,
		tagStore:         tagStore,
		userStore:        userStore,
		questionStore:    questionStore,
		interactionStore: interactionStore,
		rolandStore:      rolandStore,
		devopsStore:      devopsStore,
		kpiStore:         kpiStore,
	}, nil
}
