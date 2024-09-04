package usecase_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserUseCase_GetUserOrCreate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoErrorf(t, err, "can not open in memory db, error : %v", err)
	err = db.AutoMigrate(&entity.User{})
	assert.NoErrorf(t, err, "can not migrate in memory db, error: %v", err)
	repo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)
	mobileNumber := "090012305412"
	user, err := userUseCase.GetUserOrCreate(mobileNumber)
	assert.NoError(t, err)
	assert.Equal(t, user.MobileNumber, mobileNumber)
}
