package repository

import (
	"context"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/redis/go-redis/v9"
)

type OTPCodeRepository interface {
	Save(context.Context, *entity.OTPCode) error
	GetCode(context.Context, string) (string, error)
	DeleteCode(context.Context, string) error
}

type otpCodeRepository struct {
	client *redis.Client
}

func NewOTPCodeRepository(client *redis.Client) OTPCodeRepository {
	return &otpCodeRepository{
		client: client,
	}
}

func (otpRepo otpCodeRepository) Save(ctx context.Context, otpCode *entity.OTPCode) error {
	err := otpRepo.client.Set(ctx, otpCode.MobileNumber, otpCode.Code, time.Minute).Err()
	return err
}

func (otpRepo otpCodeRepository) GetCode(ctx context.Context, mobileNumber string) (string, error) {
	stringCmd := otpRepo.client.Get(ctx, mobileNumber)
	return stringCmd.Result()
}

func (otpRepo otpCodeRepository) DeleteCode(ctx context.Context, mobileNumber string) error {
	err := otpRepo.client.Del(ctx, mobileNumber)
	return err.Err()
}
