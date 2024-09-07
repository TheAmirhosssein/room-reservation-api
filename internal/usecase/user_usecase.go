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
	user := entity.NewUser("", mobileNumber)
	u.Repo.ByMobileNumber(mobileNumber, user)
	if user.ID == 0 {
		err := u.Repo.Save(user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (u UserUseCase) DoesUserExist(mobileNumber string) bool {
	user := entity.NewUser("", mobileNumber)
	u.Repo.ByMobileNumber(mobileNumber, user)
	return user.ID != 0
}
