package repository_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Save(t *testing.T) {
	db := database.TestDb()
	repo := repository.NewUserRepository(db)
	user := entity.NewUser("something", "09000000000")
	err := repo.Save(user)
	assert.NoErrorf(t, err, "can not save user, error: %v", err)

	var savedUser entity.User
	result := db.First(&savedUser, user.ID)
	assert.NoError(t, err, "failed to retrieve user: %v", result.Error)

	assert.Equal(t, user.MobileNumber, savedUser.MobileNumber)
	assert.Equal(t, user.FullName, savedUser.FullName)
}

func TestUserRepository_ByMobileNumber(t *testing.T) {
	db := database.TestDb()
	repo := repository.NewUserRepository(db)
	user := entity.NewUser("something", "09000000000")
	err := repo.Save(user)
	assert.NoErrorf(t, err, "can not save user, error: %v", err)

	var savedUser entity.User
	result := repo.ByMobileNumber("09000000000", &savedUser)
	assert.NoError(t, err, "failed to retrieve user: %v", result.Error)
	assert.Equal(t, user.ID, savedUser.ID)

	var wrongUser entity.User
	repo.ByMobileNumber("wrong", &wrongUser)
	assert.Zero(t, wrongUser.ID)
}
