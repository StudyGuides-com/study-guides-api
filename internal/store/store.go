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
