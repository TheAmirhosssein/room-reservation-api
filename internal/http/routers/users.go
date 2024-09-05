package routers

import (
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func UserRouters(server *gin.Engine, prefix string) {
	userRouter := server.Group(prefix)
	userRouter.POST("authenticate", handlers.Authenticate)
	userRouter.POST("token", handlers.Token)
}
