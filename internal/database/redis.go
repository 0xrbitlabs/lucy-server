package database

import (
	"context"
	"fmt"
	"os"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Println("Failed to read REDIS_URL environment variable")
		os.Exit(1)
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Printf("Error while parsing redis URL: %v", err)
		os.Exit(1)
	}
	client := redis.NewClient(opts)
	err = client.Ping(context.Background()).Err()
	if err != nil {
		fmt.Printf("Error while pinging redis: %v", err)
		os.Exit(1)
	}
  fmt.Println("Connected to Redis")
	return client
}
