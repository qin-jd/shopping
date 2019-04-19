package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreateConnection() (*gorm.DB, error) {
	host := "192.168.0.111"
	user := "mytestroot"
	dbName := "shopping"
	password := "mytestroot"
	return gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, dbName,
	),
	)
}
