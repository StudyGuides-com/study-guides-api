package tag

type tagStore struct {
}

type TagStore interface {
}


func NewTagStore() TagStore {
	return &tagStore{}
}