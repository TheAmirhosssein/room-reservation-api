package repository

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type CityRepository interface {
	Save(context.Context, *entity.City) *gorm.DB
}

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return cityRepository{db: db}
}

func (repo cityRepository) Save(ctx context.Context, city *entity.City) *gorm.DB {
	return repo.db.WithContext(ctx).Save(city)
}
