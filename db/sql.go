package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var sql *gorm.DB

func InitSql() {
	db, err := gorm.Open(mysql.Open(""))
	if err != nil {
		panic(err)
	}
	sql = db
}

func GetSql() *gorm.DB {
	return sql
}
