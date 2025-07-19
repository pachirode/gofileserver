package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/pachirode/gofileserver/internal/pkg/models"
)

type UserStore interface {
	Create(ctx context.Context, user *models.UserM) error
	Get(ctx context.Context, username string) (*models.UserM, error)
	Update(ctx context.Context, user *models.UserM) error
	List(ctx context.Context, offset, limit int) (int64, []*models.UserM, error)
	Delete(ctx context.Context, username string) error
}

type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db: db}
}

func (u *users) Create(ctx context.Context, user *models.UserM) error {
	return u.db.Create(&user).Error
}

func (u *users) Get(ctx context.Context, username string) (*models.UserM, error) {
	var user models.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *users) Update(ctx context.Context, user *models.UserM) error {
	return u.db.Save(user).Error
}

func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*models.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&models.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
