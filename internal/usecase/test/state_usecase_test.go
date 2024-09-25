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

func TestStateUseCase_Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	state := entity.NewState("something")
	err = useCase.Create(ctx, &state)
	assert.NoError(t, err)

	var savedState entity.State
	result := db.First(&savedState, 1)
	assert.NoError(t, err, "failed to retrieve State: %v", result.Error)

	assert.Equal(t, "something", savedState.Title)
}

func TestStateUseCase_GetStatesList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)

	state := entity.NewState("something")
	repo.Save(ctx, &state)
	newState := entity.NewState("something else")
	repo.Save(ctx, &newState)

	var count int64
	db.Model(&entity.State{}).Count(&count)

	states, err := useCase.GetStateList(ctx, 1, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))

	states, err = useCase.GetStateList(ctx, 1, 10, "something else")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
}

func TestStateUseCase_DoesStateExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	stateUseCase := usecase.NewStateUseCase(repo)

	result := stateUseCase.DoesStateExist(ctx, 1)
	assert.False(t, result)

	state := entity.NewState("something")
	err = repo.Save(ctx, &state).Error
	assert.NoError(t, err)
	result = stateUseCase.DoesStateExist(ctx, 1)
	assert.True(t, result)
}

func TestStateRepository_GetStateById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)

	stateUseCase := usecase.NewStateUseCase(repo)
	_, err = stateUseCase.GetStateById(ctx, 1)
	assert.Error(t, err)

	state := entity.NewState("something")
	repo.Save(ctx, &state)

	_, err = stateUseCase.GetStateById(ctx, 1)
	assert.NoError(t, err)
}

func TestStateUseCase_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)

	state := entity.NewState("something")
	err = repo.Save(ctx, &state).Error

	assert.NoError(t, err)
	err = useCase.Update(ctx, state.ID, map[string]any{"Title": "something else"})
	assert.NoError(t, err)
	repo.ById(ctx, state.ID, &state)
	assert.Equal(t, state.Title, "something else")
}

func TestStateUseCase_DeleteById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)

	state := entity.NewState("something")
	repo.Save(ctx, &state)
	var count int64
	db.Model(&entity.State{}).Count(&count)

	err = useCase.DeleteById(ctx, state.ID)
	assert.NoError(t, err)
	var countAfterDelete int64
	db.Model(&entity.State{}).Count(&countAfterDelete)

	assert.Equal(t, countAfterDelete, count-1)
}
