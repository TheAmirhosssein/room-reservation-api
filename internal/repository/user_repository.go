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
	AllUser() ([]entity.User, *gorm.DB)
	Count() (int, error)
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

func (userRepo userRepository) AllUser() ([]entity.User, *gorm.DB) {
	var users []entity.User
	query := userRepo.db.Find(&users)
	return users, query
}

func (userRepo userRepository) PaginateUsers(limit, offset int, query *gorm.DB) ([]entity.User, *gorm.DB) {
	var users []entity.User
	dbQuery := query.Limit(limit).Offset(offset).Find(&users)
	return users, dbQuery
}

func (userRepo userRepository) Count() (int, error) {
	var count int64
	err := userRepo.db.Model(&entity.User{}).Count(&count).Error
	return int(count), err
}
