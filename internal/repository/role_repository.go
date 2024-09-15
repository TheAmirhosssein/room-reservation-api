package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(*entity.Role) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return roleRepository{db: db}
}

func (roleRepo roleRepository) Save(role *entity.Role) error {
	return roleRepo.db.Save(role).Error
}
