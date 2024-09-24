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

func TestStateRepository_Test(t *testing.T) {
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
