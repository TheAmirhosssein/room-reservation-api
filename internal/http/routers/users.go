package routers

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/handlers"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouters(server *gin.Engine, prefix string) {
	userRouter := server.Group(prefix)
	userRouter.POST("authenticate", handlers.Authenticate)
	userRouter.POST("token", handlers.Token)
	userRouter.GET("me", middlewares.AuthenticateMiddleware, handlers.Me)
	userRouter.PUT("me", middlewares.AuthenticateMiddleware, handlers.UpdateUser)
	userRouter.DELETE("me", middlewares.AuthenticateMiddleware, handlers.DeleteAccount)
}
