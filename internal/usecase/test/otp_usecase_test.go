package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestOTPUseCase_GenerateCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mr, err := miniredis.Run()
	assert.NoErrorf(t, err, "An error occurred while starting miniredis: %v", err)
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	mobileNumber := "09220002200"
	otpRepo := repository.NewOTPCodeRepository(rdb)
	otpUseCase := usecase.NewOTPCase(otpRepo)
	code, err := otpUseCase.GenerateCode(ctx, mobileNumber)
	assert.NoErrorf(t, err, "An error occurred generating new otp code: %v", err)

	savedOTPCode, err := mr.Get(mobileNumber)
	assert.NoErrorf(t, err, "An error occurred getting value from miniredis: %v", err)
	assert.Equal(t, savedOTPCode, code)

	expectedError := errors.New("please wait a minute to get new code")
	_, err = otpUseCase.GenerateCode(ctx, mobileNumber)
	assert.EqualError(t, expectedError, err.Error(), "Expected error to match the expected error")
}

func TestOTPUseCase_IsCodeValid(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mr, err := miniredis.Run()
	assert.NoErrorf(t, err, "An error occurred while starting miniredis: %v", err)
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	otpRepo := repository.NewOTPCodeRepository(rdb)
	otpUseCase := usecase.NewOTPCase(otpRepo)
	mobileNumber := "09220002200"

	code, err := otpUseCase.GenerateCode(ctx, mobileNumber)
	assert.NoError(t, err)

	err = otpUseCase.ValidateCode(ctx, mobileNumber, code)
	assert.NoError(t, err)

	savedCode, _ := mr.Get(mobileNumber)
	fmt.Println(savedCode)
	err = otpUseCase.ValidateCode(ctx, mobileNumber, code)
	assert.Error(t, err)

	expectedError := errors.New("this code is invalid, please get new one")
	err = otpUseCase.ValidateCode(ctx, "wrongnumber", code)
	assert.EqualError(t, err, expectedError.Error(), "Expected error to match the expected error")

	otpUseCase.GenerateCode(ctx, mobileNumber)
	err = otpUseCase.ValidateCode(ctx, mobileNumber, "code")
	expectedError = errors.New("this code is incorrect")
	assert.EqualError(t, err, expectedError.Error(), "Expected error to match the expected error")
}
