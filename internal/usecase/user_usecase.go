package usecase

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
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

func (u UserUseCase) DoesUserExist(mobileNumber string) bool {
	user := new(entity.User)
	u.Repo.ByMobileNumber(mobileNumber, user)
	return user.ID != 0
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
