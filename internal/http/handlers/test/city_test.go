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
	addres := fmt.Sprintf("/settings/states/%v/city", state.ID)

	var countBeforeSave int64
	db.Model(&entity.City{}).Count(&countBeforeSave)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("POST", addres, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("POST", addres, bytes.NewReader(body))
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
	req, _ = http.NewRequest("POST", addres, bytes.NewReader(body))
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
	addres := fmt.Sprintf("/settings/states/%v/city", state.ID)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", addres, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", addres, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", addres, nil)
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
	addres := fmt.Sprintf("/settings/states/%v/city/1", state.ID)

	otherState, err := createState(db)
	assert.NoError(t, err)
	otherAddress := fmt.Sprintf("/settings/states/%v/city/1", otherState.ID)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", addres, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	city := entity.NewCity("something", state)
	cityRepo := repository.NewCityRepository(db)
	cityRepo.Save(context.Background(), &city)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", addres, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", addres, nil)
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
