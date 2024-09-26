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

func TestStateRepository_Save(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	state := entity.NewState("something")
	err = repo.Save(ctx, &state).Error
	assert.NoError(t, err)

	var savedState entity.State
	result := db.First(&savedState, state.ID)
	assert.NoError(t, err, "failed to retrieve State: %v", result.Error)

	assert.Equal(t, state.Title, savedState.Title)
}

func TestStateRepository_StateList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)

	state := entity.NewState("something else")
	repo.Save(ctx, &state)

	otherState := entity.NewState("something")
	repo.Save(ctx, &otherState)

	var count int64
	db.Model(&entity.State{}).Count(&count)

	states, query := repo.StateList(ctx, "")
	assert.NoError(t, query.Error)
	assert.Equal(t, int(count), len(states))

	states, query = repo.StateList(ctx, "else")
	assert.NoError(t, query.Error)
	assert.Equal(t, len(states), 1)

	states, query = repo.StateList(ctx, "something")
	assert.NoError(t, query.Error)
	assert.Equal(t, len(states), 2)
}

func TestStateRepository_Paginate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	state := entity.NewState("something else")
	repo.Save(ctx, &state)

	newState := entity.NewState("something")
	repo.Save(ctx, &newState)

	var count int64
	db.Model(&entity.State{}).Count(&count)

	_, query := repo.StateList(ctx, "")
	assert.NoError(t, query.Error)

	states, err := repo.StatePaginate(10, 0, query)
	assert.NoError(t, err)
	assert.Equal(t, len(states), 2)

	states, err = repo.StatePaginate(1, 0, query)
	assert.NoError(t, err)
	assert.Equal(t, len(states), 1)
}

func TestStateRepository_Count(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)
	state := entity.NewState("something")
	repo.Save(ctx, &state)
	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
}

func TestStateRepository_ById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)

	state := entity.State{}
	repo.ById(ctx, 1, &state)
	assert.Equal(t, uint(0), state.ID)

	state = entity.NewState("something")
	repo.Save(ctx, &state)
	assert.Equal(t, uint(1), state.ID)
}

func TestStateRepository_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	query := repo.Save(ctx, &state)
	assert.NoError(t, query.Error)

	err = repo.Update(ctx, &state, map[string]any{"title": "something else"})
	assert.NoError(t, err)
	assert.Equal(t, state.Title, "something else")
}

func TestStateRepository_Delete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewStateRepository(db)

	state := entity.NewState("something")
	repo.Save(ctx, &state)
	var count int64
	db.Model(&entity.State{}).Count(&count)

	repo.Delete(ctx, &state)
	var countAfterDelete int64
	db.Model(&entity.State{}).Count(&countAfterDelete)

	assert.Equal(t, countAfterDelete, count-1)
}
