package repository_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository_Save(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	user := entity.NewUser("something", "09000000000", entity.UserRole)
	err = repo.Save(&user)
	assert.NoErrorf(t, err, "can not save user, error: %v", err)

	var savedUser entity.User
	result := db.First(&savedUser, user.ID)
	assert.NoError(t, err, "failed to retrieve user: %v", result.Error)

	assert.Equal(t, user.MobileNumber, savedUser.MobileNumber)
	assert.Equal(t, user.FullName, savedUser.FullName)
}

func TestUserRepository_ByMobileNumber(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	user := entity.NewUser("something", "09000000000", entity.UserRole)
	err = repo.Save(&user)
	assert.NoErrorf(t, err, "can not save user, error: %v", err)

	var savedUser entity.User
	result := repo.ByMobileNumber("09000000000", &savedUser)
	assert.NoError(t, err, "failed to retrieve user: %v", result.Error)
	assert.Equal(t, user.ID, savedUser.ID)

	var wrongUser entity.User
	repo.ByMobileNumber("wrong", &wrongUser)
	assert.Zero(t, wrongUser.ID)
}

func TestUserRepository_ById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)

	user := entity.User{}
	repo.ById(1, &user)
	assert.Equal(t, uint(0), user.ID)

	user = entity.NewUser("something", "09001231021", entity.UserRole)
	repo.Save(&user)
	assert.Equal(t, uint(1), user.ID)
}

func TestUserRepository_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)

	user := entity.NewUser("something", "09001230024", entity.UserRole)
	err = repo.Save(&user)
	assert.NoError(t, err)

	err = repo.Update(&user, map[string]any{"FullName": "something else"})
	assert.NoError(t, err)
	assert.Equal(t, user.FullName, "something else")
}

func TestUserRepository_Delete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)

	user := entity.NewUser("something", "09900302020", entity.UserRole)
	repo.Save(&user)
	var count int64
	db.Model(&entity.User{}).Count(&count)

	repo.Delete(&user)
	var countAfterDelete int64
	db.Model(&entity.User{}).Count(&countAfterDelete)

	assert.Equal(t, countAfterDelete, count-1)
}
