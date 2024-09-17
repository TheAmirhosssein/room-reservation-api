package repository

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(*entity.User) error
	ByMobileNumber(string, *entity.User) *gorm.DB
	ById(uint, *entity.User) *gorm.DB
	Update(*entity.User, map[string]any) error
	Delete(*entity.User) *gorm.DB
	AllUser() ([]entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (userRepo userRepository) Save(user *entity.User) error {
	return userRepo.db.Save(user).Error
}

func (userRepo userRepository) ByMobileNumber(value string, user *entity.User) *gorm.DB {
	return userRepo.db.First(&user, "mobile_number = ?", value)
}

func (userRepo userRepository) ById(id uint, user *entity.User) *gorm.DB {
	return userRepo.db.First(&user, "ID = ?", id)
}

func (userRepo userRepository) Update(user *entity.User, newInfo map[string]any) error {
	return userRepo.db.Model(&user).Updates(newInfo).Error
}

func (userRepo userRepository) Delete(user *entity.User) *gorm.DB {
	return userRepo.db.Delete(user)
}

func (userRepo userRepository) AllUser() ([]entity.User, error) {
	var users []entity.User
	err := userRepo.db.Find(&users).Error
	return users, err
}
