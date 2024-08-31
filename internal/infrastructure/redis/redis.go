package redis

import (
	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func NewConnection(conf *config.Config) {
	opt, err := redis.ParseURL(conf.Redis.Url)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt)
}

func GetClient() *redis.Client {
	return client
}
