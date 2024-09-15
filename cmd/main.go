package main

import (
	"log"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/database"
	"github.com/TheAmirhosssein/room-reservation-api/internal/infrastructure/server"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	err = database.StartDB()
	if err != nil {
		log.Fatalf("Database error: %s", err)
	}
	server.Run(conf)
}
