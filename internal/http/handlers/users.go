package handlers

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
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
	user = *userUseCase.GetUserOrCreate(user.MobileNumber)
	response := gin.H{"message": "otp code sent", "mobile_number": user.MobileNumber}
	context.JSON(http.StatusOK, response)
}
