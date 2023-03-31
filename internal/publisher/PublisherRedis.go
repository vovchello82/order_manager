package publisher

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type PublisherRedis struct {
	redisClient *redis.Client
}

func NewRedisPublisher() *PublisherRedis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &PublisherRedis{
		redisClient: redisClient,
	}
}

func (pr *PublisherRedis) EmitObject(payload string) error {
	log.Printf("sending paload %s", payload)
	return pr.redisClient.RPush(context.TODO(), "orders", payload).Err()
}
