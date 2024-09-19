package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/routers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/redis"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createUserAndToken(userRepo repository.UserRepository, role string) (entity.User, string) {
	mobileNumber := strconv.Itoa(rand.Int())
	user := entity.NewUser("something", mobileNumber, role)
	userRepo.Save(&user)
	token, err := utils.GenerateAccessToken(user.ID, mobileNumber, user.Role)
	if err != nil {
		panic(token)
	}
	return user, token
}

func TestAuthenticateHandler(t *testing.T) {
	invalidMobileNumberResponse := `{"message":"mobile number format is not valid"}`
	invalidMobileNumber := "1234"
	body, _ := json.Marshal(map[string]string{"mobile_number": invalidMobileNumber})
	server := gin.Default()
	routers.UserRouters(server, "user")
	req, _ := http.NewRequest("POST", "/user/authenticate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, invalidMobileNumberResponse, string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	validNumber := "09001230012"
	body, _ = json.Marshal(map[string]string{"mobile_number": validNumber})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("POST", "/user/authenticate", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	invalidTimeResponse := `{"message":"please wait a minute to get new code"}`
	req, _ = http.NewRequest("POST", "/user/authenticate", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))
	redis.InitiateTestClient()
}

func TestTokenHandler(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	server := gin.Default()
	routers.UserRouters(server, "user")
	mobileNumber := "09001234141"
	code := "123456"
	body, _ := json.Marshal(map[string]string{"mobile_number": mobileNumber, "code": "wrongCode"})

	req, _ := http.NewRequest("POST", "/user/token", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	invalidTimeResponse := `{"message":"this code is invalid, please get new one"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	client := redis.TestClient()
	client.Set(context.TODO(), mobileNumber, code, time.Hour)

	invalidTimeResponse = `{"message":"this code is incorrect"}`
	req, _ = http.NewRequest("POST", "/user/token", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))

	body, _ = json.Marshal(map[string]string{"mobile_number": mobileNumber, "code": code})
	invalidTimeResponse = `{"message":"this code is incorrect"}`
	req, _ = http.NewRequest("POST", "/user/token", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	redis.InitiateTestClient()
	database.InitiateTestDB()
}

func TestMeHandler(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	user, token := createUserAndToken(userRepo, entity.UserRole)

	server := gin.Default()
	routers.UserRouters(server, "user")

	req, _ := http.NewRequest("GET", "/user/me", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", "/user/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	response, _ := io.ReadAll(w.Body)
	var result map[string]any
	json.Unmarshal(response, &result)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, user.ID, uint(result["id"].(float64)))
	assert.Equal(t, user.MobileNumber, result["mobile_number"])
}

func TestEditMeInfo(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	user, token := createUserAndToken(userRepo, entity.UserRole)

	server := gin.Default()
	routers.UserRouters(server, "user")

	body, _ := json.Marshal(map[string]string{"full_name": "something else"})
	req, _ := http.NewRequest("PUT", "/user/me", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	wrongBody, _ := json.Marshal(map[string]string{"wrongOne": "something else"})
	req, _ = http.NewRequest("PUT", "/user/me", bytes.NewBuffer(wrongBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PUT", "/user/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	userAfterUpdate := new(entity.User)
	userRepo.ById(user.ID, userAfterUpdate)
	assert.Equal(t, userAfterUpdate.FullName, "something else")
}

func TestDeleteAccount(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	user, token := createUserAndToken(userRepo, entity.UserRole)

	userRepo.Save(&user)
	var count int64
	db.Model(&entity.User{}).Count(&count)

	server := gin.Default()
	routers.UserRouters(server, "user")

	req, _ := http.NewRequest("DELETE", "/user/me", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("DELETE", "/user/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	userRepo.Delete(&user)
	var countAfterDelete int64
	db.Model(&entity.User{}).Count(&countAfterDelete)
	assert.Equal(t, countAfterDelete, count-1)

	req, _ = http.NewRequest("GET", "/user/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAllUsers(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	server := gin.Default()
	routers.UserRouters(server, "user")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", "/user/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/user/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/user/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRetrieveUser(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	server := gin.Default()
	routers.UserRouters(server, "user")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/user/users/50500", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	body, _ := json.Marshal(map[string]string{"full_name": "something", "role": "Admin"})

	server := gin.Default()
	routers.UserRouters(server, "user")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("PUT", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("PUT", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PUT", "/user/users/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updatedUser := new(entity.User)
	userRepo.ById(1, updatedUser)
	assert.Equal(t, updatedUser.FullName, "something")
	assert.Equal(t, updatedUser.Role, "Admin")

	req, _ = http.NewRequest("PUT", "/user/users/50500", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	supportBody, err := json.Marshal(map[string]string{"full_name": "something", "role": "Admin"})
	assert.NoError(t, err)
	_, supportToken := createUserAndToken(userRepo, entity.SupportRole)
	req, _ = http.NewRequest("PUT", "/user/users/1", bytes.NewBuffer(supportBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)
	fmt.Println(string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	newSupportBody, _ := json.Marshal(map[string]string{"full_name": "something else", "role": "User"})
	req, _ = http.NewRequest("PUT", "/user/users/1", bytes.NewBuffer(newSupportBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	invalidSupportBody, _ := json.Marshal(map[string]string{"full_name": "something else", "role": "ddksodksdksso"})
	req, _ = http.NewRequest("PUT", "/user/users/1", bytes.NewBuffer(invalidSupportBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUser(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	user, token := createUserAndToken(userRepo, entity.UserRole)
	adminUser, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	_, supportToken := createUserAndToken(userRepo, entity.SupportRole)

	userRepo.Save(&user)
	var count int64
	db.Model(&entity.User{}).Count(&count)

	server := gin.Default()
	routers.UserRouters(server, "user")

	req, _ := http.NewRequest("DELETE", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	req, _ = http.NewRequest("DELETE", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var countAfterDelete int64
	db.Model(&entity.User{}).Count(&countAfterDelete)
	assert.Equal(t, countAfterDelete, count-1)

	req, _ = http.NewRequest("DELETE", "/user/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	address := fmt.Sprintf("/user/users/%v", adminUser.ID)
	req, _ = http.NewRequest("DELETE", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
