package repository

import (
	"errors"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(*entity.Role) error
	ExitsByName(name string) bool
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

func (roleRepo roleRepository) ExitsByName(name string) bool {
	var role entity.Role
	err := roleRepo.db.First(&role, "name = ?", name).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
