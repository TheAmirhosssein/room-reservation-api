package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/http/handlers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateHandler(t *testing.T) {
	invalidMobileNumberResponse := `{"message":"mobile number format is not valid"}`
	invalidMobileNumber := "1234"
	body, _ := json.Marshal(map[string]string{"mobile_number": invalidMobileNumber})
	r := gin.Default()
	r.POST("/", handlers.Authenticate)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, invalidMobileNumberResponse, string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	validNumber := "09001230012"
	body, _ = json.Marshal(map[string]string{"mobile_number": validNumber})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	invalidTimeResponse := `{"message":"please wait a minute to get new code"}`
	req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))
	redis.InitiateTestClient()
}

func TestTokenHandler(t *testing.T) {
	redis.InitiateTestClient()
	database.InitiateTestDB()

	r := gin.Default()
	r.POST("/", handlers.Token)
	mobileNumber := "09001234141"
	code := "123456"
	body, _ := json.Marshal(map[string]string{"mobile_number": mobileNumber, "code": "wrongCode"})

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	invalidTimeResponse := `{"message":"this code is invalid, please get new one"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	client := redis.TestClient()
	client.Set(context.TODO(), mobileNumber, code, time.Hour)

	invalidTimeResponse = `{"message":"this code is incorrect"}`
	req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = io.ReadAll(w.Body)
	assert.Equal(t, invalidTimeResponse, string(responseData))

	body, _ = json.Marshal(map[string]string{"mobile_number": mobileNumber, "code": code})
	invalidTimeResponse = `{"message":"this code is incorrect"}`
	req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	redis.InitiateTestClient()
	database.InitiateTestDB()
}
