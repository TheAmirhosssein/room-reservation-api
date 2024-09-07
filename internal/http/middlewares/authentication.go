package middlewares

import (
	"net/http"
	"strings"

	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	if authHeader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no token have been provided"})
		return
	}
	if !strings.Contains(authHeader, "Bearer ") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "wrong token type"})
		return
	}
	token := strings.Split(authHeader, " ")[1]
	claims, err := utils.ValidateToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}
	db := database.GetDb()
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	mobileNumber := claims["mobileNumber"].(string)
	if !userUseCase.DoesUserExist(mobileNumber) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid user"})
		return
	}
	context.Set("userId", int64(claims["userId"].(float64)))
	context.Set("mobileNumber", claims["mobileNumber"].(string))
	context.Next()
}
