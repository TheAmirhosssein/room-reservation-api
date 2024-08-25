package app

import (
	"fmt"
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) {
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})
	handler.Run(fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port))
}
