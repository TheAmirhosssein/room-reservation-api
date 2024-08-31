package entity

import (
	"math/rand"
	"strconv"
)

type OTPCode struct {
	MobileNumber string
	Code         string
}

func NewOtpCode(mobileNumber string) OTPCode {
	randomNumber := rand.Intn(900000) + 100000
	return OTPCode{
		MobileNumber: mobileNumber,
		Code:         strconv.Itoa(randomNumber),
	}
}
