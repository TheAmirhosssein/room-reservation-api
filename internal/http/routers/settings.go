package routers

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/handlers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/middlewares"
	"github.com/gin-gonic/gin"
)

func SettingsRouters(server *gin.Engine, prefix string) {
	adminUser := server.Group(prefix)
	adminUser.Use(middlewares.AuthenticateMiddleware, middlewares.SupportOrAdminMiddleware)
	adminUser.POST("states", handlers.CreateState)
}
