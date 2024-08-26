package app

import (
	"fmt"
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) {
	db, err := GetDB(conf.DB.Host, conf.DB.Username, conf.DB.Username, conf.DB.DB)
	if err != nil {
		panic(err.Error())
	}
	err = Migrate(db)
	if err != nil {
		panic(err.Error())
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})
	handler.Run(fmt.Sprintf("%v:%v", conf.HTTP.Host, conf.HTTP.Port))
}
