package handlers

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/models"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
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

func StateList(context *gin.Context) {
	db := database.GetDb()
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	pageSize := utils.ParseQueryParamToInt(context.Query("page-size"), 10)
	pageNumber := utils.ParseQueryParamToInt(context.Query("page"), 1)
	title := context.Query("title")
	states, err := useCase.GetStateList(context, pageNumber, pageSize, title)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	usersCount, err := useCase.Count()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	userResponse := models.NewStateListResponse(states)
	response := utils.GenerateListResponse(userResponse, usersCount, pageSize, pageNumber)
	context.JSON(http.StatusOK, response)
}
