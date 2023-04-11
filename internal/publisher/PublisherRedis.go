package publisher

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type PublisherRedis struct {
	key         string
	redisClient *redis.Client
}

func NewRedisPublisher(key string) *PublisherRedis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &PublisherRedis{
		key:         key,
		redisClient: redisClient,
	}
}

func (pr *PublisherRedis) EmitObject(payload string) error {
	log.Printf("sending paload %s", payload)
	return pr.redisClient.RPush(context.TODO(), pr.key, payload).Err()
}
