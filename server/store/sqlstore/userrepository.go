package sqlstore

import (
	"errors"

	"github.com/Agilen/Mess/server/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	var count int64
	u.HashPassword()
	err := r.store.db.Where("name = ?", u.BaseUserInfo.Name).Find(&[]model.User{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("name is already in use")
	}

	if err := r.store.db.Create(&u).Scan(&u).Error; err != nil {
		return err
	}

	return nil
}
