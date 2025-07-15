package biz

import (
	"github.com/pachirode/gofileserver/internal/gofileserver/biz/user"
	"github.com/pachirode/gofileserver/internal/gofileserver/store"
)

type IBiz interface {
	Users() user.UserBiz
}

type biz struct {
	ds store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
