package store

import "github.com/Agilen/Mess/server/model"

type UserRepository interface {
	Create(*model.User) error
	FindUser(string, string) (bool, error)
}
