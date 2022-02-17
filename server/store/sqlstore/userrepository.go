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
	err := r.store.db.Where("name = ?", u.BaseUserInfo.Login).Find(&[]model.User{}).Count(&count).Error
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

func (r *UserRepository) FindUser(login string, password string) (bool, error) {

	user := &model.User{
		BaseUserInfo: model.BaseUserInfo{
			Login:    login,
			Password: password,
		},
	}

	user.HashPassword()

	err := r.store.db.Where(user).Find(&user).Error
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return true, nil
}
