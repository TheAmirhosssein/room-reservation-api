package routers

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/handlers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/middlewares"
	"github.com/gin-gonic/gin"
)

func SettingsRouters(server *gin.Engine, prefix string) {
	protectedRoutes := server.Group(prefix)
	protectedRoutes.Use(middlewares.AuthenticateMiddleware, middlewares.SupportOrAdminMiddleware)

	freeRoutes := server.Group(prefix)

	protectedRoutes.POST("states", handlers.CreateState)
	freeRoutes.GET("states", handlers.StateList)
	freeRoutes.GET("states/:id", handlers.RetrieveState)
	protectedRoutes.PUT("states/:id", handlers.UpdateState)
	protectedRoutes.DELETE("states/:id", handlers.DeleteState)

	protectedRoutes.POST("states/:stateId/city", handlers.CreateCity)
	freeRoutes.GET("states/:id/city", handlers.CityList)
	freeRoutes.GET("states/:id/city/:cityId", handlers.RetrieveCity)
	protectedRoutes.PUT("states/:id/city/:cityId", handlers.UpdateCity)
}
