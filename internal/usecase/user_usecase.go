package usecase

import (
	"errors"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"gorm.io/gorm"
)

type UserUseCase struct {
	Repo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return UserUseCase{Repo: userRepo}
}

func (u UserUseCase) GetUserOrCreate(mobileNumber string) (*entity.User, error) {
	user := entity.NewUser("", mobileNumber, entity.UserRole)
	u.Repo.ByMobileNumber(mobileNumber, &user)
	if user.ID == 0 {
		err := u.Repo.Save(&user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (u UserUseCase) DoesUserExist(id uint) bool {
	user := new(entity.User)
	err := u.Repo.ById(id, user).Error
	return !(errors.Is(err, gorm.ErrRecordNotFound))
}

func (u UserUseCase) GetUserById(id uint) (entity.User, error) {
	user := new(entity.User)
	query := u.Repo.ById(id, user)
	return *user, query.Error
}

func (u UserUseCase) Update(id uint, newInfo map[string]any) error {
	user, err := u.GetUserById(id)
	if err != nil {
		return err
	}
	err = u.Repo.Update(&user, newInfo)
	return err
}

func (u UserUseCase) DeleteById(id uint) error {
	user, err := u.GetUserById(id)
	if err != nil {
		return err
	}
	return u.Repo.Delete(&user).Error
}

func (u UserUseCase) AllUser(count, itemCount int) ([]entity.User, error) {
	users, query := u.Repo.UserList("", "")
	return users, query.Error
}

func (u UserUseCase) Count() (int, error) {
	return u.Repo.Count()
}
