package models

import (
	"imshit/aircraftwar/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Password []byte `gorm:"not null"`
}

func QueryUser(name string) (*User, error) {
	user := &User{}
	result := db.GetSql().Table("users").Where("name = ?", name).First(user)
	return user, result.Error
}

func CreateUser(user *User) error {
	result := db.GetSql().Table("users").Create(user)
	return result.Error
}
