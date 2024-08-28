package database

import (
	"fmt"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB(host, user, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", host, user, password, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
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
