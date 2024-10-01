package usecase

import (
	"context"
	"errors"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"gorm.io/gorm"
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

func (u CityUseCase) CityList(ctx context.Context, page, size, stateId int, title string) ([]entity.City, error) {
	_, query := u.Repo.List(ctx, title, stateId)
	if err := query.Error; err != nil {
		return nil, err
	}
	offset := utils.PageToOffset(page, size)
	state, err := u.Repo.Paginate(size, offset, query)
	return state, err
}

func (u CityUseCase) Count(ctx context.Context) (int, error) {
	return u.Repo.Count(ctx)
}

func (u CityUseCase) DoesCityExist(ctx context.Context, id, stateId uint) bool {
	city := new(entity.City)
	err := u.Repo.ById(ctx, id, city).Error
	return !(errors.Is(err, gorm.ErrRecordNotFound))
}

func (u CityUseCase) ById(ctx context.Context, id uint) (entity.City, error) {
	city := new(entity.City)
	query := u.Repo.ById(ctx, id, city)
	return *city, query.Error
}

func (u CityUseCase) Update(ctx context.Context, id uint, newInfo map[string]any) error {
	city, err := u.ById(ctx, id)
	if err != nil {
		return err
	}
	err = u.Repo.Update(ctx, &city, newInfo)
	return err
}

func (u CityUseCase) DeleteById(ctx context.Context, id uint) error {
	city, err := u.ById(ctx, id)
	if err != nil {
		return err
	}
	return u.Repo.Delete(ctx, &city).Error
}
