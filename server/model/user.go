package model

import (
	"crypto/sha256"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LastActive time.Time `gorm:"not null"`
	BaseUserInfo
}

type BaseUserInfo struct {
	gorm.Model
	Name     string `gorm:"not null;unique;"`
	Email    string `gorm:"not null;type:varchar(100)"`
	Password string `gorm:"not null;type:varchar(100)"`
}

func (u *User) HashPassword() {
	sum := sha256.Sum256([]byte(u.BaseUserInfo.Password))
	u.BaseUserInfo.Password = string(sum[:])
}
