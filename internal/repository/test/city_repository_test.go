package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCityRepository_Save(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	stateRepo := repository.NewStateRepository(db)
	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	cityRepo := repository.NewCityRepository(db)
	city := entity.NewCity("something", state)
	err = cityRepo.Save(ctx, &city).Error
	assert.NoError(t, err)

	var savedCity entity.City
	result := db.Preload("State").First(&savedCity, city.ID)
	assert.NoError(t, err, "failed to retrieve State: %v", result.Error)

	assert.Equal(t, city.Title, savedCity.Title)
	assert.Equal(t, savedCity.State.ID, state.ID)
}
