package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

func QueryUser(name string) (*User, error) {
	// TODO
	return &User{}, nil
}

func AddUser(user *User) error {
	return nil
}
