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
	Count(context.Context) (int, error)
	ById(context.Context, uint, *entity.State) *gorm.DB
	Update(context.Context, *entity.State, map[string]any) error
	Delete(context.Context, *entity.State) *gorm.DB
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

func (repo stateRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := repo.db.Model(&entity.State{}).Count(&count).Error
	return int(count), err
}

func (repo stateRepository) ById(ctx context.Context, id uint, state *entity.State) *gorm.DB {
	return repo.db.WithContext(ctx).First(&state, "ID = ?", id)
}

func (repo stateRepository) Update(ctx context.Context, state *entity.State, newInfo map[string]any) error {
	return repo.db.WithContext(ctx).Model(&state).Updates(newInfo).Error
}

func (repo stateRepository) Delete(ctx context.Context, state *entity.State) *gorm.DB {
	return repo.db.WithContext(ctx).Delete(state)
}
