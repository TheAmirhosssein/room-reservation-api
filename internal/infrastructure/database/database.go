package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{})
}

func StartDB() {
	db := GetDb()
	err := migrate(db)
	if err != nil {
		panic(err.Error())
	}
}

func TestDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	return db
}

func GetDb() *gorm.DB {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return TestDb()
		}
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
