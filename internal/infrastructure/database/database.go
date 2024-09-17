package database

import (
	"fmt"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(host, user, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", host, user, password, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{})
}

func StartDB() error {
	db := GetDb()
	err := Migrate(db)
	if err != nil {
		return err
	}
	return err
}

func GetDb() *gorm.DB {
	if config.InTestMode() {
		return TestDb()
	}

	conf, err := config.NewConfig()
	if err != nil {
		panic(err.Error())
	}
	db, err := initDB(conf.DB.Host, conf.DB.Username, conf.DB.Username, conf.DB.DB)
	if err != nil {
		panic(err.Error())
	}
	return db
}
