package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheAmirhosssein/room-reservation-api/internal/http/models"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateCity(context *gin.Context) {
	stateId, err := strconv.ParseInt(context.Param("stateId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	db := database.GetDb()
	stateRepo := repository.NewStateRepository(db)
	stateUseCase := usecase.NewStateUseCase(stateRepo)
	if !stateUseCase.DoesStateExist(context, uint(stateId)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "state not found"})
		return
	}
	state, err := stateUseCase.GetStateById(context, uint(stateId))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	body := new(models.City)
	err = context.BindJSON(body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	cityRepo := repository.NewCityRepository(db)
	cityUseCase := usecase.NewCityUseCase(cityRepo)
	city, err := cityUseCase.Create(context, body.Title, state)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	response := models.NewCityResponse(city)
	context.JSON(http.StatusCreated, response)
}
