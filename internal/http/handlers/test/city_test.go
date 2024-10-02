package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/routers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/redis"
	"github.com/TheAmirhosssein/room-reservation-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createState(db *gorm.DB) (entity.State, error) {
	stateRepo := repository.NewStateRepository(db)
	state := entity.NewState("something")
	return state, stateRepo.Save(context.Background(), &state).Error
}

func TestCreateCity(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	body, err := json.Marshal(map[string]string{"title": "something"})
	assert.NoError(t, err)

	state, err := createState(db)
	assert.NoError(t, err)
	address := fmt.Sprintf("/settings/states/%v/city", state.ID)

	var countBeforeSave int64
	db.Model(&entity.City{}).Count(&countBeforeSave)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("POST", address, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("POST", address, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var countAfterSave int64
	db.Model(&entity.City{}).Count(&countAfterSave)
	assert.Equal(t, countBeforeSave+1, countAfterSave)

	city := new(entity.City)
	cityRepo := repository.NewCityRepository(db)
	cityRepo.ById(context.Background(), uint(1), city)
	assert.Equal(t, city.State.ID, state.ID)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("POST", address, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest("POST", "/settings/states/505050/city", bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListCity(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	state, err := createState(db)
	assert.NoError(t, err)
	address := fmt.Sprintf("/settings/states/%v/city", state.ID)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/settings/states/505050/city", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRetrieveCity(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	state, err := createState(db)
	assert.NoError(t, err)
	address := fmt.Sprintf("/settings/states/%v/city/1", state.ID)

	otherState, err := createState(db)
	assert.NoError(t, err)
	otherAddress := fmt.Sprintf("/settings/states/%v/city/1", otherState.ID)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	city := entity.NewCity("something", state)
	cityRepo := repository.NewCityRepository(db)
	cityRepo.Save(context.Background(), &city)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/settings/states/505050/city", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", otherAddress, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCity(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	state, err := createState(db)
	assert.NoError(t, err)
	address := fmt.Sprintf("/settings/states/%v/city/1", state.ID)
	otherAddress := fmt.Sprintf("/settings/states/%v/city/3", state.ID)

	cityRepo := repository.NewCityRepository(db)
	newCity := entity.NewCity("na", state)
	err = cityRepo.Save(context.Background(), &newCity).Error
	assert.NoError(t, err)

	body, _ := json.Marshal(map[string]string{"title": "something"})

	server := gin.Default()
	routers.SettingsRouters(server, "settings")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("PUT", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("PUT", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PUT", address, bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updatedCity := new(entity.City)
	cityRepo.ById(req.Context(), 1, updatedCity)
	assert.Equal(t, updatedCity.Title, "something")

	req, _ = http.NewRequest("PUT", otherAddress, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.NoError(t, err)
	_, supportToken := createUserAndToken(userRepo, entity.SupportRole)

	supportBody, _ := json.Marshal(map[string]string{"title": "something"})
	req, _ = http.NewRequest("PUT", address, bytes.NewBuffer(supportBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	invalidSupportBody, _ := json.Marshal(map[string]string{"titledfdfdf": "something"})
	req, _ = http.NewRequest("PUT", address, bytes.NewBuffer(invalidSupportBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCity(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	user, token := createUserAndToken(userRepo, entity.UserRole)
	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	_, supportToken := createUserAndToken(userRepo, entity.SupportRole)

	state, err := createState(db)
	assert.NoError(t, err)
	address := fmt.Sprintf("/settings/states/%v/city", state.ID)

	cityRepo := repository.NewCityRepository(db)
	newCity := entity.NewCity("title", state)
	err = cityRepo.Save(context.Background(), &newCity).Error
	assert.NoError(t, err)

	userRepo.Save(&user)
	var count int64
	db.Model(&entity.City{}).Count(&count)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")

	req, _ := http.NewRequest("DELETE", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	req, _ = http.NewRequest("DELETE", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var countAfterDelete int64
	db.Model(&entity.City{}).Count(&countAfterDelete)
	assert.Equal(t, countAfterDelete, count-1)

	req, _ = http.NewRequest("DELETE", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("DELETE", address, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
