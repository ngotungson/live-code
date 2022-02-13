package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var db *gorm.DB
var once sync.Once

func Init() error {
	once.Do(func() {
		er := initDbInstance()
		if er != nil {
			panic(er)
		}
	})

	return nil
}

func GetDB() *gorm.DB {
	once.Do(func() {
		er := initDbInstance()
		if er != nil {
			panic(er)
		}
	})
	return db
}

func initDbInstance() error {
	connStr := "user=postgres password=098poiA# port=5432 dbname=testdb sslmode=disable"
	connection, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	db = connection
	return nil
}
