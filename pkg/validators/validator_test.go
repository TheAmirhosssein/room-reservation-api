package validators_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/pkg/validators"
)

func TestValidateMobileNumber(t *testing.T) {
	validNumbers := []string{
		"09123456789",
		"989123456789",
		"+989123456789",
		"00989123456789",
	}

	for _, num := range validNumbers {
		if !validators.ValidateMobileNumber(num) {
			t.Errorf("Expected valid mobile number, but got invalid for: %s", num)
		}
	}

	invalidNumbers := []string{
		"1234567890",
		"0912345678",
		"00981234567890",
		"abcdefghijk",
		"09123456789abc",
	}

	for _, num := range invalidNumbers {
		if validators.ValidateMobileNumber(num) {
			t.Errorf("Expected invalid mobile number, but got valid for: %s", num)
		}
	}
}
