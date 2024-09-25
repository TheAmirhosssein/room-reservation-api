package repository

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type StateRepository interface {
	Save(context.Context, *entity.State) *gorm.DB
	StateList(context.Context, string) ([]entity.State, *gorm.DB)
	Paginate(int, int, *gorm.DB) ([]entity.State, error)
	Count() (int, error)
}

type stateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return stateRepository{db: db}
}

func (repo stateRepository) Save(ctx context.Context, state *entity.State) *gorm.DB {
	return repo.db.WithContext(ctx).Save(state)
}

func (repo stateRepository) StateList(ctx context.Context, title string) ([]entity.State, *gorm.DB) {
	var states []entity.State
	query := repo.db.WithContext(ctx).Model(&entity.State{}).
		Where("title LIKE ?", "%"+title+"%").
		Find(&states)
	return states, query
}

func (repo stateRepository) Paginate(limit, offset int, query *gorm.DB) ([]entity.State, error) {
	var states []entity.State
	err := query.Limit(limit).Offset(offset).Find(&states).Error
	return states, err
}

func (repo stateRepository) Count() (int, error) {
	var count int64
	err := repo.db.Model(&entity.State{}).Count(&count).Error
	return int(count), err
}
