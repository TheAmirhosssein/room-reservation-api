package redis

import (
	"os"
	"strings"

	"github.com/TheAmirhosssein/room-reservation-api/config"
	"github.com/alicebob/miniredis/v2"
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

func TestClient() *redis.Client {
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
			return TestClient()
		}
	}
	return Client()
}
