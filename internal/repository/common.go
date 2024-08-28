package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type commonBehavior[T entity.DBTable] struct {
	db *gorm.DB
}

func NewCommonBehavior[T entity.DBTable](db *gorm.DB) BaseRepository[T] {
	return &commonBehavior[T]{
		db: db,
	}
}

func (c *commonBehavior[T]) GetDB() *gorm.DB {
	return c.db
}

func (c *commonBehavior[T]) Save(model *T) error {
	return c.db.Save(model).Error
}
