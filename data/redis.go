package data

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)



func RedisSetup() *redis.Client {
	_ = godotenv.Load()
	connStr := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(connStr)

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return client
  }