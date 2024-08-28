package server

import (
	"fmt"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/http/routers"
	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) {
	server := gin.Default()
	routers.UserRouters(server, "user/")
	server.Run(fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port))
}
