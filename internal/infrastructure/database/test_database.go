package database

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDb *gorm.DB

func InitiateTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	testDb = db
}

func TestDb() *gorm.DB {
	if testDb == nil {
		InitiateTestDB()
	}
	return testDb
}
