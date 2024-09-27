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
	Count(context.Context) (int, error)
	ById(context.Context, uint, *entity.City) *gorm.DB
	Update(context.Context, *entity.City, map[string]any) error
	Delete(context.Context, *entity.City) *gorm.DB
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

func (repo cityRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := repo.db.Model(&entity.City{}).Count(&count).Error
	return int(count), err
}

func (repo cityRepository) ById(ctx context.Context, id uint, city *entity.City) *gorm.DB {
	return repo.db.WithContext(ctx).Preload("State").First(&city, "ID = ?", id)
}

func (repo cityRepository) Update(ctx context.Context, city *entity.City, newInfo map[string]any) error {
	return repo.db.WithContext(ctx).Model(&city).Updates(newInfo).Error
}

func (repo cityRepository) Delete(ctx context.Context, city *entity.City) *gorm.DB {
	return repo.db.WithContext(ctx).Delete(city)
}
