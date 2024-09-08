package middlewares_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/middlewares"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	database.InitiateTestDB()
	db := database.GetDb()

	server := gin.Default()
	server.GET("/", middlewares.AuthenticateMiddleware, func(ctx *gin.Context) {
		mobileNumber := ctx.GetString("mobileNumber")
		userId := ctx.GetUint("userId")
		ctx.JSON(http.StatusOK, gin.H{"id": userId, "mobile_number": mobileNumber})
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ := io.ReadAll(w.Body)
	expectedResponse := `{"message":"no token have been provided"}`
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Equal(t, string(response), expectedResponse)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "something")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ = io.ReadAll(w.Body)
	expectedResponse = `{"message":"wrong token type"}`
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Equal(t, string(response), expectedResponse)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer something")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ = io.ReadAll(w.Body)
	expectedResponse = `{"message":"invalid token"}`
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Equal(t, string(response), expectedResponse)

	mobileNumber := "09001110011"
	token, err := utils.GenerateAccessToken(1, mobileNumber)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ = io.ReadAll(w.Body)
	expectedResponse = `{"message":"invalid user"}`
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Equal(t, string(response), expectedResponse)

	user := entity.NewUser("something", mobileNumber)
	userRepo := repository.NewUserRepository(db)
	userRepo.Save(&user)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ = io.ReadAll(w.Body)
	expectedResponse = fmt.Sprintf(`{"id":%v,"mobile_number":"%v"}`, user.ID, user.MobileNumber)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, expectedResponse, string(response))
}
