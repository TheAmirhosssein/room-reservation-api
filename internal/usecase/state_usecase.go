package usecase

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
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

func (u StateUseCase) Count() (int, error) {
	return u.Repo.Count()
}
