package db

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var sql *gorm.DB

func InitSql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
		viper.GetString("mysql.timeout"),
	)
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	sql = _db
}

func GetSql() *gorm.DB {
	return sql
}
