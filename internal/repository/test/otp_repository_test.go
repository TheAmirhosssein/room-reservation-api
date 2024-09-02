package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/alicebob/miniredis/v2"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestOTPCodeRepository_Save(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("An error occurred while starting miniredis: %v", err)
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	mobileNumber := "09000000000"
	otpRepo := repository.NewOTPCodeRepository(rdb)

	otpCode := entity.NewOtpCode(mobileNumber)

	err = otpRepo.Save(ctx, &otpCode)
	assert.NoError(t, err, "Save should not return an error")

	savedCode, err := mr.Get(otpCode.MobileNumber)
	assert.NoError(t, err, "Save should not return an error")
	assert.Equal(t, otpCode.Code, savedCode, "The saved code should match the input")
}

func TestOTPCodeRepository_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("An error occurred while starting miniredis: %v", err)
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	mobileNumber := "09000000000"
	otpRepo := repository.NewOTPCodeRepository(rdb)

	otpCode := entity.NewOtpCode(mobileNumber)

	err = otpRepo.Save(ctx, &otpCode)
	assert.NoError(t, err, "Save should not return an error")

	savedCode, err := mr.Get(otpCode.MobileNumber)
	assert.NoError(t, err, "Save should not return an error")
	code, err := otpRepo.GetCode(ctx, mobileNumber)
	assert.NoError(t, err, "Get Code should not return an error")
	assert.Equal(t, code, savedCode, "The saved code should match the input")
}
