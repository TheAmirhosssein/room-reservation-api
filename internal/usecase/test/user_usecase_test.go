package usecase_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
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
	database.Migrate(db)
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
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)

	MobileNumber := "09001230101"
	result := userUseCase.DoesUserExist(1)
	assert.False(t, result)

	user := entity.NewUser("something", MobileNumber, entity.UserRole)
	err = repo.Save(&user)
	assert.NoError(t, err)
	result = userUseCase.DoesUserExist(1)
	assert.True(t, result)
}

func TestUserRepository_GetUserById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)

	userUseCase := usecase.NewUserUseCase(repo)
	_, err = userUseCase.GetUserById(1)
	assert.Error(t, err)

	user := entity.NewUser("something", "09001231010", entity.UserRole)
	repo.Save(&user)

	_, err = userUseCase.GetUserById(1)
	assert.NoError(t, err)
}

func TestUserUseCase_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)

	user := entity.NewUser("something", "09001234565", entity.UserRole)
	err = repo.Save(&user)

	assert.NoError(t, err)
	err = useCase.Update(user.ID, map[string]any{"FullName": "something else"})
	assert.NoError(t, err)
	repo.ById(user.ID, &user)
	assert.Equal(t, user.FullName, "something else")
}

func TestUserUseCase_DeleteById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)

	user := entity.NewUser("something", "09900302020", entity.UserRole)
	repo.Save(&user)
	var count int64
	db.Model(&entity.User{}).Count(&count)

	err = useCase.DeleteById(user.ID)
	assert.NoError(t, err)
	var countAfterDelete int64
	db.Model(&entity.User{}).Count(&countAfterDelete)

	assert.Equal(t, countAfterDelete, count-1)
}

func TestUserUseCase_GetUsersList(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)

	user := entity.NewUser("something", "09900302020", entity.UserRole)
	repo.Save(&user)
	newUser := entity.NewUser("something else", "09900302023", entity.UserRole)
	repo.Save(&newUser)

	var count int64
	db.Model(&entity.User{}).Count(&count)

	users, err := useCase.GetUsersList(1, 1, "", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))

	users, err = useCase.GetUsersList(1, 10, "09900302023", "something else")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
}

func TestUserUseCase_Count(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.Migrate(db)
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)
	user := entity.NewUser("something", "09900302020", entity.UserRole)
	repo.Save(&user)
	count, err := useCase.Count()
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
}
