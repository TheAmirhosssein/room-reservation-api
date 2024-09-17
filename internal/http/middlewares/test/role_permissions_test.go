package middlewares_test

import (
	"fmt"
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

func TestAdminRoleMiddleware(t *testing.T) {
	database.InitiateTestDB()
	db := database.GetDb()

	server := gin.Default()
	server.GET("/", middlewares.AuthenticateMiddleware, middlewares.AdminRoleMiddleware, func(ctx *gin.Context) {
		mobileNumber := ctx.GetString("mobileNumber")
		userId := ctx.GetUint("userId")
		ctx.JSON(http.StatusOK, gin.H{"id": userId, "mobile_number": mobileNumber})
	})

	userRepo := repository.NewUserRepository(db)

	user := entity.NewUser("user", "09002520066", entity.UserRole)
	err := userRepo.Save(&user)
	assert.NoError(t, err)
	userToken, err := utils.GenerateAccessToken(user.ID, user.MobileNumber, user.Role)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusForbidden)

	adminUser := entity.NewUser("admin", "09002520023", entity.AdminRole)
	err = userRepo.Save(&adminUser)
	assert.NoError(t, err)
	adminToken, err := utils.GenerateAccessToken(adminUser.ID, adminUser.MobileNumber, entity.AdminRole)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestSupportRoleMiddleware(t *testing.T) {
	database.InitiateTestDB()
	db := database.GetDb()

	server := gin.Default()
	server.GET("/", middlewares.AuthenticateMiddleware, middlewares.SupportRoleMiddleware, func(ctx *gin.Context) {
		mobileNumber := ctx.GetString("mobileNumber")
		userId := ctx.GetUint("userId")
		ctx.JSON(http.StatusOK, gin.H{"id": userId, "mobile_number": mobileNumber})
	})

	userRepo := repository.NewUserRepository(db)

	user := entity.NewUser("user", "09002520066", entity.UserRole)
	err := userRepo.Save(&user)
	assert.NoError(t, err)
	userToken, err := utils.GenerateAccessToken(user.ID, user.MobileNumber, user.Role)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userToken))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusForbidden)

	adminUser := entity.NewUser("admin", "09002520023", entity.AdminRole)
	err = userRepo.Save(&adminUser)
	assert.NoError(t, err)
	adminToken, err := utils.GenerateAccessToken(adminUser.ID, adminUser.MobileNumber, user.Role)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", adminToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusForbidden)

	supportUser := entity.NewUser("support", "09002520023", entity.SupportRole)
	err = userRepo.Save(&adminUser)
	assert.NoError(t, err)
	supportToken, err := utils.GenerateAccessToken(supportUser.ID, supportUser.MobileNumber, supportUser.Role)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", supportToken))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusForbidden)
}
