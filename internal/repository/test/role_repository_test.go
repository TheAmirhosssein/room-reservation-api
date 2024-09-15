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

func TestRoleRepository_Test(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = database.Migrate(db)
	assert.NoError(t, err)
	var countBeForeSave int64
	db.Model(&entity.Role{}).Count(&countBeForeSave)
	assert.Zero(t, countBeForeSave)

	role := entity.NewRole("something")
	repo := repository.NewRoleRepository(db)
	err = repo.Save(&role)
	assert.NoError(t, err)

	var countAfterSave int64
	db.Model(&entity.Role{}).Count(&countAfterSave)
	assert.Equal(t, countAfterSave, int64(1))
}

func TestRoleRepository_ExitsByName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = database.Migrate(db)
	assert.NoError(t, err)
	repo := repository.NewRoleRepository(db)

	exist := repo.ExitsByName("something")
	assert.False(t, exist)

	role := entity.NewRole("something")
	err = repo.Save(&role)
	assert.NoError(t, err)

	exist = repo.ExitsByName("something")
	assert.True(t, exist)
}
