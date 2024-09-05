package redis

import (
	"os"
	"strings"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/alicebob/miniredis/v2"
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

func GetTestClient() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return rdb
}

func GetClient() *redis.Client {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return GetTestClient()
		}
	}
	return client
}
