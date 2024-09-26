package repository

import (
	"gorm.io/gorm"
)

type CityRepository interface {
	// Save(context.Context, *entity.State) *gorm.DB
}

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return cityRepository{db: db}
}
