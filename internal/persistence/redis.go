package persistence

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/tanmancan/draw-together/internal/config"
)

var redisClient *redis.Client

func init() {
	h := config.AppConfig.NetworkServiceRedisClient.Host
	p := config.AppConfig.NetworkServiceRedisClient.Port
	addr := fmt.Sprintf("%s:%d", h, p)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func GetClient() *redis.Client {
	return redisClient
}
