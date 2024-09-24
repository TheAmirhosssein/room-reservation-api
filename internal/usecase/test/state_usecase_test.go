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
	err = useCase.Create(ctx, "something")
	assert.NoError(t, err)

	var savedState entity.State
	result := db.First(&savedState, 1)
	assert.NoError(t, err, "failed to retrieve State: %v", result.Error)

	assert.Equal(t, "something", savedState.Title)
}
