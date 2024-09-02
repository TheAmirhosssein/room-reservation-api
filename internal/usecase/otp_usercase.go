package usecase

import (
	"errors"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
)

type OTPUseCase struct {
	Repo repository.OTPCodeRepository
}

func NewOTPCase(otpRepo repository.OTPCodeRepository) OTPUseCase {
	return OTPUseCase{Repo: otpRepo}
}

func (otp OTPUseCase) GenerateCode(mobileNumber string) error {
	if otp.Repo.GetCode(mobileNumber) != "" {
		return errors.New("please wait a minute to get new code")
	}
	otpCode := entity.NewOtpCode(mobileNumber)
	err := otp.Repo.Save(&otpCode)
	return err
}
