package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func SetupRedis() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	Redis = redis.NewClient(&redis.Options{Addr: host + ":" + port, Password: os.Getenv("REDIS_PASSWORD"), DB: 0})
}
