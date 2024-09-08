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
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	repo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)
	mobileNumber := "090012305412"
	var count int64
	db.Model(&entity.User{}).Count(&count)
	user, err := userUseCase.GetUserOrCreate(mobileNumber)
	assert.NoError(t, err)
	assert.Equal(t, user.MobileNumber, mobileNumber)
	var countAfter int64
	db.Model(&entity.User{}).Count(&countAfter)
	count += 1
	assert.Equal(t, count, countAfter)
	_, err = userUseCase.GetUserOrCreate(mobileNumber)
	assert.NoError(t, err)
	assert.Equal(t, count, countAfter)
}

func TestUserUseCase_DoesUserExist(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	repo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)

	MobileNumber := "09001230101"
	result := userUseCase.DoesUserExist(MobileNumber)
	assert.False(t, result)

	user := entity.NewUser("something", MobileNumber)
	err = repo.Save(&user)
	assert.NoError(t, err)
	result = userUseCase.DoesUserExist(MobileNumber)
	assert.True(t, result)
}

func TestUserRepository_GetUserById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	repo := repository.NewUserRepository(db)

	userUseCase := usecase.NewUserUseCase(repo)
	_, err = userUseCase.GetUserById(1)
	assert.Error(t, err)

	user := entity.NewUser("something", "09001231010")
	repo.Save(&user)

	_, err = userUseCase.GetUserById(1)
	assert.NoError(t, err)
}

func TestUserUseCase_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{})
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)

	user := entity.NewUser("something", "")
	err = useCase.Update(&user)
	assert.Error(t, err)
	assert.EqualError(t, err, "mobile number can not be empty")

	user = entity.NewUser("something", "09001234565")
	err = useCase.Update(&user)
	assert.Error(t, err)
	assert.NotEqual(t, err.Error(), "mobile number can not be empty")

	err = repo.Save(&user)
	updateUser := entity.NewUser("something else", "09001234565")
	assert.NoError(t, err)
	err = useCase.Update(&updateUser)
	assert.NoError(t, err)
	repo.ById(user.ID, &user)
	assert.Equal(t, user.ID, updateUser.ID)
	assert.Equal(t, user.FullName, updateUser.FullName)
}
