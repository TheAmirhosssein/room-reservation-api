package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
)

type BaseRepository[T entity.DBTable] interface {
	Save(*T) error
}

type UserRepository interface {
	BaseRepository[entity.User]
}
