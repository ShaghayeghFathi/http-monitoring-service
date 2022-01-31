package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Urls     []Url  `gorm:"foreignkey:user_id"`
}

func NewUser(username, password string) (*User, error) {
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}
	pass, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{Username: username, Password: pass}, nil
}

func HashPassword(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

func (user *User) ValidatePassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) == nil
}
