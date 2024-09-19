package validators_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/pkg/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateMobileNumber(t *testing.T) {
	validNumbers := []string{
		"09123456789",
		"989123456789",
		"+989123456789",
		"00989123456789",
	}

	for _, num := range validNumbers {
		assert.True(t, validators.ValidateMobileNumber(num), "Expected valid mobile number, but got invalid for: %s", num)
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

func TestIsRoleValid(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{role: "Admin", expected: true},
		{role: "Support", expected: true},
		{role: "User", expected: true},
		{role: "UnknownRole", expected: false},
		{role: "admin", expected: false},
	}

	for _, test := range tests {
		t.Run(test.role, func(t *testing.T) {
			result := validators.IsRoleValid(test.role)
			if result != test.expected {
				t.Errorf("For role %s, expected %v but got %v", test.role, test.expected, result)
			}
		})
	}
}
