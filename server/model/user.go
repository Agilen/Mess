package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string
	Password   string
	LastActive time.Time
}
