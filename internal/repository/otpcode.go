package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/redis/go-redis/v9"
)

type otpCodeRepository struct {
	client *redis.Client
}

func NewOTPCodeRepository(client *redis.Client) OTPCodeRepository {
	return &otpCodeRepository{
		client: client,
	}
}

func (otpRepo otpCodeRepository) Save(otpCode *entity.OTPCode) {
	err := otpRepo.client.Set(context.TODO(), otpCode.MobileNumber, otpCode.Code, time.Minute).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (otpRepo otpCodeRepository) GetCode(mobileNumber string) string {
	stringCmd := otpRepo.client.Get(context.TODO(), mobileNumber)
	return stringCmd.Val()
}

func (otpRepo otpCodeRepository) GenerateCode(otpCode *entity.OTPCode) error {
	if otpRepo.GetCode(otpCode.MobileNumber) != "" {
		return errors.New("please wait a minute to get new code")
	}
	otpRepo.Save(otpCode)
	return nil
}
