/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

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
