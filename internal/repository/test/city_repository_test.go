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

func TestCityRepository_List(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	stateRepo := repository.NewStateRepository(db)
	repo := repository.NewCityRepository(db)

	state := entity.NewState("something else")
	stateRepo.Save(ctx, &state)

	city := entity.NewCity("something else", state)
	repo.Save(ctx, &city)

	otherCity := entity.NewCity("something", state)
	repo.Save(ctx, &otherCity)

	var count int64
	db.Model(&entity.City{}).Count(&count)

	cities, query := repo.List(ctx, "", 1)
	assert.NoError(t, query.Error)
	assert.Equal(t, int(count), len(cities))

	cities, query = repo.List(ctx, "else", 1)
	assert.NoError(t, query.Error)
	assert.Equal(t, len(cities), 1)

	cities, query = repo.List(ctx, "", 0)
	assert.NoError(t, query.Error)
	assert.Equal(t, len(cities), 2)
}

func TestCityRepository_Paginate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.Migrate(db)
	repo := repository.NewCityRepository(db)

	state := entity.NewState("something")
	stateRepo := repository.NewStateRepository(db)
	stateRepo.Save(ctx, &state)

	city := entity.NewCity("something else", state)
	repo.Save(ctx, &city)

	newCity := entity.NewCity("something", state)
	repo.Save(ctx, &newCity)

	var count int64
	db.Model(&entity.City{}).Count(&count)

	_, query := repo.List(ctx, "", 1)
	assert.NoError(t, query.Error)

	cities, err := repo.Paginate(10, 0, query)
	assert.NoError(t, err)
	assert.Equal(t, len(cities), 2)

	cities, err = repo.Paginate(1, 0, query)
	assert.NoError(t, err)
	assert.Equal(t, len(cities), 1)
}
