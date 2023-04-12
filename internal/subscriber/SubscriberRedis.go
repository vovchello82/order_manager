package subscriber

import (
	"context"
	"encoding/json"
	"log"
	domain "order_manager/internal/db"

	"github.com/redis/go-redis/v9"
)

type SubscriberRedis struct {
	topic       string
	redisClient *redis.Client
	db          *domain.Storage
}

func NewSubscriberRedis(topic string, storage *domain.Storage) *SubscriberRedis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &SubscriberRedis{
		topic:       topic,
		db:          storage,
		redisClient: redisClient,
	}
}

func (sr *SubscriberRedis) HandleOrderResult() error {
	topic := sr.redisClient.Subscribe(context.Background(), sr.topic)
	defer topic.Close()

	log.Printf("start listening on %s", sr.topic)

	channel := topic.Channel()

	for msg := range channel {
		result := &OrderPlacementResult{}
		json.Unmarshal([]byte(msg.Payload), result)
		log.Printf("order was placed %s by taxi %s", result.OrderId, result.TaxiId)
		sr.db.DeleteOrderById(result.OrderId)
	}

	return nil
}
