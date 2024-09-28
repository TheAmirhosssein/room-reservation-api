package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCityUseCase_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.Migrate(db)

	cityRepo := repository.NewCityRepository(db)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	cityUseCase := usecase.NewCityUseCase(cityRepo)

	citiesCount, err := cityRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, citiesCount, 0)

	city, err := cityUseCase.Create(ctx, "something", state)
	assert.NoError(t, err)
	assert.Equal(t, city.Title, "something")

	countAfterCreate, err := cityRepo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, countAfterCreate, citiesCount+1)
}
