package utils

import (
	"errors"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() (string, error) {
	if config.InTestMode() {
		return "secretestkey", nil
	}

	conf, err := config.NewConfig()
	if err != nil {
		return "", err
	}

	return conf.APP.SecretKey, nil
}

func GenerateAccessToken(userId int64, mobileNumber string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":       userId,
		"mobileNumber": mobileNumber,
		"exp":          time.Now().Add(time.Hour * (24 * 365)).Unix(),
	})
	secretKey, err := getSecretKey()
	if err != nil {
		return "", err
	}
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(token string) (map[string]any, error) {
	paredToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid type")
		}
		secretKey, err := getSecretKey()
		if err != nil {
			return nil, err
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !paredToken.Valid {
		return nil, errors.New("invalid token")
	}
	return paredToken.Claims.(jwt.MapClaims), nil
}
