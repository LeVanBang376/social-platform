package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	log.Println("Connected to Redis")

	return client
}
