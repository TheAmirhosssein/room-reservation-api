package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type usersRepository struct {
	BaseRepository[entity.User]
}

func NewUsersRepository(db *gorm.DB) UserRepository {
	return &usersRepository{
		BaseRepository: NewCommonBehavior[entity.User](db),
	}
}
