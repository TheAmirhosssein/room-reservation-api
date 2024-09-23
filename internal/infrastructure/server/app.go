package server

import (
	"fmt"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/docs"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/routers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(conf *config.Config) {
	server := gin.Default()
	routers.UserRouters(server, "/api/v1/users")
	docs.SwaggerInfo.BasePath = "/api/v1"
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	server.Run(fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port))
}
