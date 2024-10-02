package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

// Authenticate godoc
// @Summary Authenticate User
// @Description Authenticates a user by their mobile number and generates an OTP code.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param authenticateUser body models.Authenticate true "User authentication data"
// @Success 200 {object} map[string]interface{} "OTP code sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid mobile number format or other errors"
// @Router /user/authenticate [post]
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

// Token godoc
// @Summary Validate OTP and Generate Token
// @Description Validates the OTP for a given mobile number and generates an access token.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param token body models.Token true "OTP validation data"
// @Success 200 {object} map[string]interface{} "Access token generated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid OTP or mobile number"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/token [post]
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
	accessToken, err := utils.GenerateAccessToken(user.ID, user.MobileNumber, user.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": accessToken})
}

// Me godoc
// @Summary Get User Information
// @Description Retrieves the authenticated user's details.
// @Tags User
// @Produce  json
// @Success 200 {object} models.UserResponse "User details retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/me [get]
// @Security BearerAuth
func Me(context *gin.Context) {
	db := database.GetDb()
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userId := context.GetUint("userId")
	user, err := userUseCase.GetUserById(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	userResponse := models.NewUserResponse(user)
	context.JSON(http.StatusOK, userResponse)
}

// UpdateUser godoc
// @Summary Update User Information
// @Description Updates the authenticated user's details such as their full name.
// @Tags User
// @Accept  json
// @Produce  json
// @Param updateUser body models.UpdateUser true "User update data"
// @Success 200 {object} models.UserResponse "User details updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/me [put]
// @Security BearerAuth
func UpdateUser(context *gin.Context) {
	body := new(models.UpdateUser)
	err := context.BindJSON(body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	db := database.GetDb()
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)
	data := map[string]any{"FullName": body.FullName}
	id := context.GetUint("userId")
	err = useCase.Update(id, data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	user, err := useCase.GetUserById(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	userResponse := models.NewUserResponse(user)
	context.JSON(http.StatusOK, userResponse)
}

// DeleteAccount godoc
// @Summary Delete User Account
// @Description Deletes the authenticated user's account.
// @Tags User
// @Produce  json
// @Success 204 "User account deleted successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/me [delete]
// @Security BearerAuth
func DeleteAccount(context *gin.Context) {
	db := database.GetDb()
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)
	err := useCase.DeleteById(context.GetUint("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

// AllUsers godoc
// @Summary Retrieve All Users
// @Description Fetches a paginated list of users, with optional filters for mobile number and full name.
// @Tags User
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page-size query int false "Number of users per page" default(10)
// @Param mobile-number query string false "Filter by mobile number"
// @Param full-name query string false "Filter by full name"
// @Success 200 {object} utils.PaginatedResponse{result=[]models.UserResponse} "List of users retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/users [get]
// @Security BearerAuth
func AllUsers(context *gin.Context) {
	db := database.GetDb()
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)
	pageSize := utils.ParseQueryParamToInt(context.Query("page-size"), 10)
	pageNumber := utils.ParseQueryParamToInt(context.Query("page"), 1)
	mobileNumber := context.Query("mobile-number")
	fullName := context.Query("full-name")
	allUser, err := useCase.GetUsersList(pageNumber, pageSize, mobileNumber, fullName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	usersCount, err := useCase.Count()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	userResponse := models.NewUserListResponse(allUser)
	response := utils.GenerateListResponse(userResponse, usersCount, pageSize, pageNumber)
	context.JSON(http.StatusOK, response)
}

// RetrieveUser godoc
// @Summary Retrieve User Information
// @Description Fetches details of a user by their ID.
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse "User details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/users/{id} [get]
// @Security BearerAuth
func RetrieveUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid endpoint"})
		return
	}
	db := database.GetDb()
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	user, err := userUseCase.GetUserById(uint(id))
	if !userUseCase.DoesUserExist(user.ID) {
		context.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	userResponse := models.NewUserResponse(user)
	context.JSON(http.StatusOK, userResponse)
}

// EditUser godoc
// @Summary Edit User Information
// @Description Updates user details by ID, including full name and role.
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param updateUser body models.AdminUpdateUser true "User update data"
// @Success 200 {object} models.UserResponse "User details updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or role"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/users/{id} [put]
// @Security BearerAuth
func EditUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid endpoint"})
		return
	}
	db := database.GetDb()
	repo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)
	if !userUseCase.DoesUserExist(uint(id)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	var updateData map[string]any
	role := context.GetString("role")
	data := new(models.AdminUpdateUser)
	err = context.BindJSON(data)
	if err != nil {
		context.JSONP(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !validators.IsRoleValid(data.Role) {
		context.JSONP(http.StatusBadRequest, gin.H{"message": "invalid role"})
		return
	}
	if role == entity.SupportRole && data.Role == entity.AdminRole {
		context.JSONP(http.StatusBadRequest, gin.H{"message": "invalid role to select"})
		return
	}
	updateData = map[string]any{"full_name": data.FullName, "role": data.Role}

	err = userUseCase.Update(uint(id), updateData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	user, err := userUseCase.GetUserById(uint(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	response := models.NewUserResponse(user)
	context.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Deletes a user by their ID.
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 403 {object} map[string]interface{} "Forbidden: insufficient permissions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/users/{id} [delete]
// @Security BearerAuth
func DeleteUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid endpoint"})
		return
	}
	db := database.GetDb()
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	if !userUseCase.DoesUserExist(uint(id)) {
		context.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	deleteUser, err := userUseCase.GetUserById(uint(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong!"})
		return
	}
	userRole := context.GetString("role")
	if userRole == entity.SupportRole && deleteUser.Role == entity.AdminRole {
		context.JSON(http.StatusForbidden, gin.H{"message": "you have no permission to perform this action"})
		return
	}
	userUseCase.DeleteById(uint(id))
	context.JSON(http.StatusNoContent, nil)
}
