package usecase

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
)

type OTPUseCase struct {
	Repo repository.OTPCodeRepository
}

func NewOTPCase(otpRepo repository.OTPCodeRepository) OTPUseCase {
	return OTPUseCase{Repo: otpRepo}
}

func (otp OTPUseCase) GenerateCode(ctx context.Context, mobileNumber string) (string, error) {
	code, err := otp.Repo.GetCode(ctx, mobileNumber)
	if err != nil && err != redis.Nil {
		return "", err
	}
	if code != "" {
		return "", errors.New("please wait a minute to get new code")
	}
	otpCode := entity.NewOtpCode(mobileNumber)
	err = otp.Repo.Save(ctx, &otpCode)
	if err != nil {
		return "", err
	}
	return otpCode.Code, nil
}
