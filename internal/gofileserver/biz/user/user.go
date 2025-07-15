package user

import "github.com/pachirode/gofileserver/internal/gofileserver/store"

type UserBiz interface{}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}
