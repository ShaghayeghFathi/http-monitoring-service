package db

import (
	"fmt"

	"github.com/ShaghayeghFathi/http-monitoring-service/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Setup(databaseName string) *gorm.DB {
	db := newDB(databaseName)
	migrate(db)
	db.LogMode(true)
	return db

}

func newDB(name string) *gorm.DB {

	db, err := gorm.Open("sqlite3", "./"+name)
	if err != nil {
		fmt.Println("Error in creating database file : ", err)
		return nil
	}
	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{}, &model.Request{}, &model.Url{})
}
