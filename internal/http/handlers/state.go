package handlers

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/models"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func CreateState(context *gin.Context) {
	body := new(models.State)
	err := context.BindJSON(body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	state := entity.NewState(body.Title)
	db := database.GetDb()
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	err = useCase.Create(context, &state)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	response := models.NewStateResponse(state)
	context.JSON(http.StatusCreated, response)

}
