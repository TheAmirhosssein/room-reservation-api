package usecase

import (
	"context"
	"errors"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"gorm.io/gorm"
)

type StateUseCase struct {
	Repo repository.StateRepository
}

func NewStateUseCase(repo repository.StateRepository) StateUseCase {
	return StateUseCase{Repo: repo}
}

func (u StateUseCase) Create(context context.Context, state *entity.State) error {
	return u.Repo.Save(context, state).Error
}

func (u StateUseCase) GetStateList(ctx context.Context, page, pageSize int, title string) ([]entity.State, error) {
	_, query := u.Repo.StateList(ctx, title)
	if err := query.Error; err != nil {
		return nil, err
	}
	offset := utils.PageToOffset(page, pageSize)
	state, err := u.Repo.Paginate(pageSize, offset, query)
	return state, err
}

func (u StateUseCase) Count(ctx context.Context) (int, error) {
	return u.Repo.Count(ctx)
}

func (u StateUseCase) DoesStateExist(ctx context.Context, id uint) bool {
	state := new(entity.State)
	err := u.Repo.ById(ctx, id, state).Error
	return !(errors.Is(err, gorm.ErrRecordNotFound))
}

func (u StateUseCase) GetStateById(ctx context.Context, id uint) (entity.State, error) {
	state := new(entity.State)
	query := u.Repo.ById(ctx, id, state)
	return *state, query.Error
}

func (u StateUseCase) Update(ctx context.Context, id uint, newInfo map[string]any) error {
	state, err := u.GetStateById(ctx, id)
	if err != nil {
		return err
	}
	err = u.Repo.Update(ctx, &state, newInfo)
	return err
}

func (u StateUseCase) DeleteById(ctx context.Context, id uint) error {
	state, err := u.GetStateById(ctx, id)
	if err != nil {
		return err
	}
	return u.Repo.Delete(ctx, &state).Error
}
