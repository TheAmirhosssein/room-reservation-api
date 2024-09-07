package utils

import (
	"log"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * (24 * 365)).Unix(),
	})

	if config.InTestMode() {
		return token.SignedString([]byte("secretestkey"))
	}

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	secretKey := conf.APP.SecretKey
	return token.SignedString([]byte(secretKey))
}
