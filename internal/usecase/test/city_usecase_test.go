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

func TestCityUseCase_List(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewCityRepository(db)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	useCase := usecase.NewCityUseCase(repo)

	city := entity.NewCity("something", state)
	repo.Save(ctx, &city)
	newCity := entity.NewCity("something else", state)
	repo.Save(ctx, &newCity)

	var count int64
	db.Model(&entity.City{}).Count(&count)

	cities, err := useCase.CityList(ctx, 1, 1, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cities))

	cities, err = useCase.CityList(ctx, 1, 10, 1, "something else")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cities))
}

func TestCityUseCase_DoesCityExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewCityRepository(db)
	cityUseCase := usecase.NewCityUseCase(repo)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	result := cityUseCase.DoesCityExist(ctx, 1, 1)
	assert.False(t, result)

	otherSate := entity.NewState("something")
	err = stateRepo.Save(ctx, &otherSate).Error
	assert.NoError(t, err)

	result = cityUseCase.DoesCityExist(ctx, 1, otherSate.ID)
	assert.False(t, result)

	city := entity.NewCity("something", state)
	err = repo.Save(ctx, &city).Error
	assert.NoError(t, err)
	result = cityUseCase.DoesCityExist(ctx, 1, 1)
	assert.True(t, result)
}

func TestCityRepository_GetCityById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewCityRepository(db)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	cityUseCase := usecase.NewCityUseCase(repo)
	_, err = cityUseCase.ById(ctx, 1)
	assert.Error(t, err)

	city := entity.NewCity("something", state)
	repo.Save(ctx, &city)

	_, err = cityUseCase.ById(ctx, 1)
	assert.NoError(t, err)
}

func TestCityUseCase_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewCityRepository(db)
	useCase := usecase.NewCityUseCase(repo)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	city := entity.NewCity("something", state)
	err = repo.Save(ctx, &city).Error

	assert.NoError(t, err)
	city, err = useCase.Update(ctx, city.ID, map[string]any{"Title": "something else"})
	assert.NoError(t, err)
	repo.ById(ctx, city.ID, &city)
	assert.Equal(t, city.Title, "something else")
}

func TestCityUseCase_DeleteById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.Migrate(db)
	repo := repository.NewCityRepository(db)
	useCase := usecase.NewCityUseCase(repo)
	stateRepo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	err = stateRepo.Save(ctx, &state).Error
	assert.NoError(t, err)

	city := entity.NewCity("something", state)
	repo.Save(ctx, &city)
	var count int64
	db.Model(&entity.City{}).Count(&count)

	err = useCase.DeleteById(ctx, city.ID)
	assert.NoError(t, err)
	var countAfterDelete int64
	db.Model(&entity.City{}).Count(&countAfterDelete)

	assert.Equal(t, countAfterDelete, count-1)
}
