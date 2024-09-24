package repository

import (
	"gorm.io/gorm"
)

type StateRepository interface {
	// Save(entity.State) *gorm.DB
}

type stateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return stateRepository{db: db}
}
