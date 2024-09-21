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

func AllUsers(context *gin.Context) {
	db := database.GetDb()
	repo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(repo)
	pageSize := utils.ParseQueryParamToInt(context.Query("page-size"), 10)
	pageNumber := utils.ParseQueryParamToInt(context.Query("page"), 1)
	fmt.Println(pageNumber, pageSize)
	allUser, err := useCase.AllUser(pageNumber, pageSize)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	response := models.NewUserListResponse(allUser)
	context.JSON(http.StatusOK, response)
}

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
