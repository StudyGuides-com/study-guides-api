package store

import (
	"github.com/studyguides-com/study-guides-api/internal/store/search"
	"github.com/studyguides-com/study-guides-api/internal/store/tag"
	"github.com/studyguides-com/study-guides-api/internal/store/user"
)


type Store interface {
	SearchStore() search.SearchStore
	TagStore() tag.TagStore
	UserStore() user.UserStore
}

type store struct {
	searchStore search.SearchStore
	tagStore    tag.TagStore
	userStore   user.UserStore
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

func NewStore() Store {
	return &store{
		searchStore: search.NewAlgoliaSearchStore(),
		tagStore:    tag.NewTagStore(),
		userStore:   user.NewUserStore(),
	}
}