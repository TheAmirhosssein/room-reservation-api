package app

import (
	"fmt"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(host, user, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", host, user, password, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{})
}
