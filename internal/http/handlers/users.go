package handlers

import (
	"fmt"
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/models"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/redis"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/validators"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	authenticateUser := models.Authenticate{}
	err := context.BindJSON(&authenticateUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if !validators.ValidateMobileNumber(authenticateUser.MobileNumber) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "mobile number format is not valid"})
		return
	}
	userRepo := repository.NewUserRepository(database.GetDb())
	userUseCase := usecase.NewUserUseCase(userRepo)
	user, err := userUseCase.GetUserOrCreate(authenticateUser.MobileNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	otpRepo := repository.NewOTPCodeRepository(redis.GetClient())
	otpUseCase := usecase.NewOTPCase(otpRepo)
	code, err := otpUseCase.GenerateCode(context, user.MobileNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Printf("code was %v\n", code)
	response := gin.H{"message": "otp code sent", "mobile_number": user.MobileNumber}
	context.JSON(http.StatusOK, response)
}

func Token(context *gin.Context) {
	body := models.Token{}
	err := context.BindJSON(&body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	otpRepo := repository.NewOTPCodeRepository(redis.GetClient())
	otpUseCase := usecase.NewOTPCase(otpRepo)
	err = otpUseCase.ValidateCode(context, body.MobileNumber, body.Code)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userRepo := repository.NewUserRepository(database.GetDb())
	userUseCase := usecase.NewUserUseCase(userRepo)
	var user entity.User
	userUseCase.Repo.ByMobileNumber(body.MobileNumber, &user)
	accessToken, err := utils.GenerateAccessToken(int64(user.ID), user.MobileNumber)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": accessToken})
}
