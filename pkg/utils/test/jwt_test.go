package utils_test

import (
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createTokenWithDifferentSigningMethod() string {
	secretKey := "secretestkey"
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user": "testUser",
	})
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func createInvalidToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      int64(12),
		"mobileNumber": "mobileNumber",
		"exp":          time.Now(),
	})
	secretKey := "secretestkey"
	return token.SignedString([]byte(secretKey))
}

func TestTokenValidation(t *testing.T) {
	tokenWithDifferentSigningMethod := createTokenWithDifferentSigningMethod()
	_, err := utils.ValidateToken(tokenWithDifferentSigningMethod)
	assert.Error(t, err)

	invalidToken, err := createInvalidToken()
	assert.NoError(t, err)

	_, err = utils.ValidateToken(invalidToken)
	assert.Error(t, err)

	_, err = utils.ValidateToken("invalidToken")
	assert.Error(t, err)

	_, err = utils.GenerateAccessToken(1, "something", entity.UserRole)
	assert.NoError(t, err)
}
