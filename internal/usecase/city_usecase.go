package usecase

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
)

type CityUseCase struct {
	Repo repository.CityRepository
}

func NewCityUseCase(repo repository.CityRepository) CityUseCase {
	return CityUseCase{Repo: repo}
}

func (u CityUseCase) Create(context context.Context, title string, state entity.State) (entity.City, error) {
	city := entity.NewCity(title, state)
	err := u.Repo.Save(context, &city).Error
	return city, err
}
