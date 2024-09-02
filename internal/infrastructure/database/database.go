package database

import (
	"fmt"
	"sync"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	DB   *gorm.DB
)

func initDB(host, user, password, name string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", host, user, password, name)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			db = nil
			return
		}
		DB = db
	})
	return db, err
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{})
}

func StartDB(conf *config.Config) {
	db, err := initDB(conf.DB.Host, conf.DB.Username, conf.DB.Username, conf.DB.DB)
	if err != nil {
		panic(err.Error())
	}
	err = migrate(db)
	if err != nil {
		panic(err.Error())
	}
}

func GetDb() *gorm.DB {
	return DB
}
