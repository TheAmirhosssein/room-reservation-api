package repository

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type CityRepository interface {
	Save(context.Context, *entity.City) *gorm.DB
	List(context.Context, string, int) ([]entity.City, *gorm.DB)
	Paginate(int, int, *gorm.DB) ([]entity.City, error)
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

func (repo cityRepository) List(ctx context.Context, title string, stateId int) ([]entity.City, *gorm.DB) {
	var cities []entity.City
	query := repo.db.WithContext(ctx).Preload("State").Model(&entity.City{}).Find(&cities)
	if stateId != 0 {
		query.Where("title LIKE ? AND state_id = ?", "%"+title+"%", stateId).Find(&cities)
	}
	return cities, query
}

func (repo cityRepository) Paginate(limit, offset int, query *gorm.DB) ([]entity.City, error) {
	var cities []entity.City
	err := query.Limit(limit).Offset(offset).Find(&cities).Error
	return cities, err
}
