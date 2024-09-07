package redis

import (
	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/redis/go-redis/v9"
)

func Client() *redis.Client {
	conf, err := config.NewConfig()
	if err != nil {
		panic(err.Error())
	}
	opt, err := redis.ParseURL(conf.Redis.Url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func GetClient() *redis.Client {
	if config.InTestMode() {
		return TestClient()
	}
	return Client()
}
