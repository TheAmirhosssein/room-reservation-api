package main

import (
	"log"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/app/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/app/server"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	database.StartDB(conf)
	server.Run(conf)
}
