package data

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)



func RedisSetup() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(connStr)

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return client
  }