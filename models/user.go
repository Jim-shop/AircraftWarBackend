package models

import (
	"imshit/aircraftwar/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func QueryUser(name string) (*User, error) {
	user := &User{}
	db.GetSql().Table("users").Where("name", name).First(user)

	// TODO
	return user, nil
}

func CreateUser(user *User) error {
	db.GetSql().Table("users").Create(user)
	return nil
}
