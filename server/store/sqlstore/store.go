package sqlstore

import (
	"github.com/Agilen/Mess/server/store"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	db             *gorm.DB
	UserRepository *UserRepository
}

func New(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func NewDB(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	if databaseURL[:4] == "user" {
		db, err = gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return nil, err

		}
	} else if databaseURL[:4] == "host" {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	} else {
		db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}
