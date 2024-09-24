package repository

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type StateRepository interface {
	Save(context.Context, *entity.State) *gorm.DB
}

type stateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return stateRepository{db: db}
}

func (stateRepo stateRepository) Save(ctx context.Context, state *entity.State) *gorm.DB {
	return stateRepo.db.WithContext(ctx).Save(state)
}
