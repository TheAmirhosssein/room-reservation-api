package entity_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestOTPCodeEntity_New(t *testing.T) {
	mobileNumber := "09000000000"
	otpCode := entity.NewOtpCode(mobileNumber)
	assert.Equal(t, 6, len(otpCode.Code), "code length must be 6")
	newOtpCode := entity.NewOtpCode(mobileNumber)
	assert.NotEqual(t, newOtpCode.Code, otpCode.Code, "code must be random")
}
