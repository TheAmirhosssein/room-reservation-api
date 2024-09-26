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
)

func TestCreateState(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)
	body, err := json.Marshal(map[string]string{"title": "something"})
	assert.NoError(t, err)

	var countBeforeSave int64
	db.Model(&entity.State{}).Count(&countBeforeSave)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")
	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("POST", "/settings/states", bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("POST", "/settings/states", bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var countAfterSave int64
	db.Model(&entity.State{}).Count(&countAfterSave)
	assert.Equal(t, countBeforeSave+1, countAfterSave)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("POST", "/settings/states", bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestStateList(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", "/settings/states", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/settings/states", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/settings/states", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRetrieveState(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	db := database.TestDb()
	userRepo := repository.NewUserRepository(db)

	stateRepo := repository.NewStateRepository(db)
	newState := entity.NewState("something")
	err := stateRepo.Save(context.Background(), &newState).Error
	assert.NoError(t, err)

	server := gin.Default()
	routers.SettingsRouters(server, "settings")

	_, userToken := createUserAndToken(userRepo, entity.UserRole)
	req, _ := http.NewRequest("GET", "/settings/states/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, adminToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/settings/states/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/settings/states/50500", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	_, supportToken := createUserAndToken(userRepo, entity.AdminRole)
	req, _ = http.NewRequest("GET", "/settings/states/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
