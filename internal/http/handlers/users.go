package handlers

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/app/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	user := entity.User{}
	err := context.BindJSON(&user)
	if err != nil {
		context.JSONP(http.StatusBadRequest, err.Error())
		return
	}
	userRepo := repository.NewUsersRepository(database.DB)
	userRepo.GetUserOrCreate(user.MobileNumber, &user)
	// todo: add otp code
	response := gin.H{"message": "otp code sent", "mobile_number": user.MobileNumber}
	context.JSONP(http.StatusOK, response)
}
