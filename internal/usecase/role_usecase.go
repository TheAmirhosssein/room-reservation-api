package usecase

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
)

type RoleUseCase struct {
	repo repository.RoleRepository
}

func NewRoleUserCase(repo repository.RoleRepository) RoleUseCase {
	return RoleUseCase{repo: repo}
}

func (usecase RoleUseCase) SetUpRoles(rolesName []string) error {
	for _, roleName := range rolesName {
		if !usecase.repo.ExitsByName(roleName) {
			role := entity.NewRole(roleName)
			err := usecase.repo.Save(&role)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
