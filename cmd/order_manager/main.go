package main

import (
	"fmt"
	"log"
	domain "order_manager/internal/db"
	"order_manager/internal/placeorder"
	"order_manager/internal/publisher"
	"order_manager/internal/subscriber"
	"os"
	"os/signal"
	"time"
)

func cleanup() {
	fmt.Println("cleanup")
}

const ORDER_QUEUE = "orders"
const ORDER_RESULT_TOPIC = "order-results"

func main() {
	// Hello world, the web server

	log.Println("Listing for requests at http://localhost:8000/hello")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	orderStorage := domain.NewStorage()
	redisPublisher := publisher.NewRedisPublisher(ORDER_QUEUE)
	redisSubscriber := subscriber.NewSubscriberRedis(ORDER_RESULT_TOPIC, orderStorage)

	go redisSubscriber.HandleOrderResult()

	placeOrder := placeorder.NewPlaceOrderUseCase(redisPublisher, orderStorage)

	for {
		order := domain.NewRandomOrder()
		log.Printf("sending order %s", order)
		if err := placeOrder.PlaceOrder(*order); err != nil {
			log.Fatalf("error on sending order %s", err)
		}

		time.Sleep(2 * time.Second)
	}
}
