package usecase_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRoleUseCase_SetupRoles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.Migrate(db)
	repo := repository.NewRoleRepository(db)
	useCase := usecase.NewRoleUserCase(repo)

	var countBeForeSave int64
	db.Model(&entity.Role{}).Count(&countBeForeSave)
	assert.Zero(t, countBeForeSave)

	roles := []string{"something", "something1"}
	err = useCase.SetUpRoles(roles)
	assert.NoError(t, err)

	var countAfterSave int64
	db.Model(&entity.Role{}).Count(&countAfterSave)
	assert.Equal(t, countAfterSave, int64(len(roles)))

	err = useCase.SetUpRoles(roles)
	assert.NoError(t, err)
	assert.Equal(t, countAfterSave, int64(len(roles)))
}
