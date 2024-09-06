package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var testClient *redis.Client

func InitiateClient() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	testClient = rdb
}

func TestClient() *redis.Client {
	if testClient == nil {
		InitiateClient()
	}
	return testClient
}
