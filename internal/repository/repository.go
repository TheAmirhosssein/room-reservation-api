package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type BaseRepository[T entity.DBTable] interface {
	GetDB() *gorm.DB
	Save(*T) error
}

type UserRepository interface {
	BaseRepository[entity.User]
	ByMobileNumber(string, *entity.User)
}
