package models

import (
	"imshit/aircraftwar/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null;size: 256"`
	Password []byte `gorm:"not null;size: 256"`
}

func QueryUser(name string) (*User, error) {
	user := &User{}
	result := db.GetSql().Table("users").Where("name = ?", name).First(user)
	return user, result.Error
}

func GetUser(id int) (*User, error) {
	user := &User{}
	result := db.GetSql().Table("users").Where("ID = ?", id).First(user)
	return user, result.Error
}

func CreateUser(user *User) error {
	result := db.GetSql().Table("users").Create(user)
	return result.Error
}
