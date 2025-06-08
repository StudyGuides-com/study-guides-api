package user

type userStore struct {
}

type UserStore interface {
}


func NewUserStore() UserStore {
	return &userStore{}
}