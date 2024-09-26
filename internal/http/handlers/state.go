package handlers

import (
	"net/http"
	"strconv"

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
	usersCount, err := useCase.Count(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	userResponse := models.NewStateListResponse(states)
	response := utils.GenerateListResponse(userResponse, usersCount, pageSize, pageNumber)
	context.JSON(http.StatusOK, response)
}

func RetrieveState(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	db := database.GetDb()
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	if !useCase.DoesStateExist(context, uint(id)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "state not found"})
		return
	}
	state, err := useCase.GetStateById(context, uint(id))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	response := models.NewStateResponse(state)
	context.JSON(http.StatusOK, response)
}

func UpdateState(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	db := database.GetDb()
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	if !useCase.DoesStateExist(context, uint(id)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "state not found"})
		return
	}
	body := new(models.State)
	err = context.BindJSON(body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = useCase.Update(context, uint(id), map[string]any{"title": body.Title})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	state, err := useCase.GetStateById(context, uint(id))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	response := models.NewStateResponse(state)
	context.JSON(http.StatusOK, response)
}

func DeleteState(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	db := database.GetDb()
	repo := repository.NewStateRepository(db)
	useCase := usecase.NewStateUseCase(repo)
	if !useCase.DoesStateExist(context, uint(id)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "state not found"})
		return
	}
	err = useCase.DeleteById(context, uint(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusNoContent, nil)
}
