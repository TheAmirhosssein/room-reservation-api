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
	userRepo.ByMobileNumber(user.MobileNumber, &user)
	if user.ID == 0 {
		userRepo.Save(&user)
	}
	context.JSONP(http.StatusOK, user)
}
