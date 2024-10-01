package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheAmirhosssein/room-reservation-api/internal/http/models"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/internal/usecase"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
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

func CityList(context *gin.Context) {
	stateId, err := strconv.ParseInt(context.Param("id"), 10, 64)
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
	pageSize := utils.ParseQueryParamToInt(context.Query("page-size"), 10)
	pageNumber := utils.ParseQueryParamToInt(context.Query("page"), 1)
	title := context.Query("title")
	cityRepo := repository.NewCityRepository(db)
	cityUseCase := usecase.NewCityUseCase(cityRepo)
	cities, err := cityUseCase.CityList(context, pageNumber, pageSize, int(stateId), title)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	citiesCount, err := cityUseCase.Count(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	city_list := models.NewCityListResponse(cities)
	response := utils.GenerateListResponse(city_list, citiesCount, pageSize, pageNumber)
	context.JSON(http.StatusOK, response)
}

func RetrieveCity(context *gin.Context) {
	stateId, err := strconv.ParseInt(context.Param("id"), 10, 64)
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
	cityId, err := strconv.ParseInt(context.Param("cityId"), 10, 64)
	cityRepo := repository.NewCityRepository(db)
	cityUseCase := usecase.NewCityUseCase(cityRepo)
	if !cityUseCase.DoesCityExist(context, uint(cityId), uint(stateId)) {
		context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	city, err := cityUseCase.ById(context, uint(cityId))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	response := models.NewCityResponse(city)
	context.JSON(http.StatusOK, response)
}
