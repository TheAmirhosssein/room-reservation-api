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

// CreateState handles the creation of a new state.
//
// @Summary      Create a new state
// @Description  This endpoint creates a new state record in the database.
// @Tags         states
// @Accept       json
// @Produce      json
// @Param        state  body      models.State         true  "State data"
// @Success      201    {object}  models.StateResponse  "Created state"
// @Failure      400    {object}  map[string]string     "Bad request"
// @Failure      500    {object}  map[string]string     "Internal server error"
// @Router       /settings/states [post]
// @Security BearerAuth
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

// StateList retrieves a list of states with optional filtering and pagination.
//
// @Summary      Get list of states
// @Description  This endpoint retrieves a paginated list of states. You can filter the results by title.
// @Tags         states
// @Accept       json
// @Produce      json
// @Param        page        query     int    false  "Page number"         default(1)
// @Param        page-size   query     int    false  "Page size"           default(10)
// @Param        title       query     string false  "Filter by state title"
// @Success      200         {object}  utils.PaginatedResponse{result=[]models.StateResponse}  "List of states"
// @Failure      500         {object}  map[string]string     "Internal server error"
// @Router       /settings/states [get]
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

// RetrieveState retrieves a specific state by its ID.
//
// @Summary      Get state by ID
// @Description  This endpoint retrieves the details of a specific state by its ID.
// @Tags         states
// @Accept       json
// @Produce      json
// @Param        id   path      int   true   "State ID"
// @Success      200  {object}  models.StateResponse  "State details"
// @Failure      400  {object}  map[string]string     "Invalid state ID"
// @Failure      404  {object}  map[string]string     "State not found"
// @Failure      500  {object}  map[string]string     "Failed to retrieve state"
// @Router       /settings/states/{id} [get]
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

// UpdateState updates a specific state by its ID.
//
// @Summary      Update state by ID
// @Description  This endpoint updates the details of a specific state by its ID.
// @Tags         states
// @Accept       json
// @Produce      json
// @Param        id    path      int            true   "State ID"
// @Param        body  body      models.State   true   "State data to update"
// @Success      200   {object}  models.StateResponse  "Updated state"
// @Failure      400   {object}  map[string]string     "Invalid request"
// @Failure      404   {object}  map[string]string     "State not found"
// @Failure      500   {object}  map[string]string     "Failed to update state"
// @Router       /settings/states/{id} [put]
// @Security BearerAuth
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

// DeleteState deletes a specific state by its ID.
//
// @Summary      Delete state by ID
// @Description  This endpoint deletes a specific state from the database using its ID.
// @Tags         states
// @Accept       json
// @Produce      json
// @Param        id   path      int   true   "State ID"
// @Success      204  "State deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid state ID"
// @Failure      404  {object}  map[string]string  "State not found"
// @Failure      500  {object}  map[string]string  "Failed to delete state"
// @Router       /settings/states/{id} [delete]
// @Security BearerAuth
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
