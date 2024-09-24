package usecase

import (
	"context"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
)

type StateUseCase struct {
	Repo repository.StateRepository
}

func NewStateUseCase(repo repository.StateRepository) StateUseCase {
	return StateUseCase{Repo: repo}
}

func (u StateUseCase) Create(context context.Context, title string) error {
	state := entity.NewState(title)
	return u.Repo.Save(context, &state).Error
}
