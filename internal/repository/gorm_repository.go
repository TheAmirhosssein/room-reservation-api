package repository

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type GormRepository interface {
	Save(context.Context, *entity.Model) *gorm.DB
}

type gormRepository[Model entity.Model] struct {
	db *gorm.DB
}

func NewGormRepository[Model entity.Model](db *gorm.DB) GormRepository {
	return gormRepository[Model]{db: db}
}

func (repo gormRepository[Model]) Save(ctx context.Context, object *entity.Model) *gorm.DB {
	return repo.db.WithContext(ctx).Save(object)
}
