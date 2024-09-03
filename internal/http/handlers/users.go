package handlers

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/redis"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/validators"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	user := entity.User{}
	err := context.BindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if !validators.ValidateMobileNumber(user.MobileNumber) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "mobile number format is not valid"})
		return
	}
	userRepo := repository.NewUserRepository(database.GetDb())
	userUseCase := usecase.NewUserUseCase(userRepo)
	saveUser, err := userUseCase.GetUserOrCreate(user.MobileNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	otpRepo := repository.NewOTPCodeRepository(redis.GetClient())
	otpUseCase := usecase.NewOTPCase(otpRepo)
	err = otpUseCase.GenerateCode(context, saveUser.MobileNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	response := gin.H{"message": "otp code sent", "mobile_number": saveUser.MobileNumber}
	context.JSON(http.StatusOK, response)
}
