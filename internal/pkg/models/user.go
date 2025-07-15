package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/pachirode/gofileserver/pkg/auth"
)

type UserM struct {
	Id            int64     `gorm:"column:id;primary_key"`
	Username      string    `gorm:"column:username;not null"`
	Password      string    `gorm:"column:password;not null"`
	Nickname      string    `gorm:"column:nickname"`
	Email         string    `gorm:"column:email"`
	Status        string    `gorm:"column:status"`
	Lastlogintime time.Time `gorm:"column:lastLogintime"`
	Lastip        string    `gorm:"column:lastip"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdateAt      time.Time `gorm:"column:updateAt"`
}

func (u *UserM) TableName() string {
	return "tb_http_user"
}

func (u *UserM) EncryptPassword(dbg *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
