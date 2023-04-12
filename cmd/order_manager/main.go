package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	domain "order_manager/internal/db"
	"order_manager/internal/placeorder"
	"order_manager/internal/publisher"
	"order_manager/internal/subscriber"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

var ticker = time.NewTicker(5 * time.Second)

const ORDER_QUEUE = "orders"
const ORDER_RESULT_TOPIC = "order-results"

func cleanup() {
	fmt.Println("cleanup")
	ticker.Stop()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	go serveMetrics()

	orderStorage := domain.NewStorage()
	redisPublisher := publisher.NewRedisPublisher(ORDER_QUEUE)
	redisSubscriber := subscriber.NewSubscriberRedis(ORDER_RESULT_TOPIC, orderStorage)

	go redisSubscriber.HandleOrderResult()

	placeOrder := placeorder.NewPlaceOrderUseCase(redisPublisher, orderStorage)

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("order_manager")
	gauge, err := meter.Int64ObservableGauge("order_pressure", instrument.WithDescription("orders pressure"))
	if err != nil {
		log.Fatal(err)
	}

	meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		o.ObserveInt64(gauge, int64(orderStorage.Size()))
		return nil
	}, gauge)

	for range ticker.C {
		order := domain.NewRandomOrder()
		log.Printf("sending order %s", order)
		if err := placeOrder.PlaceOrder(*order); err != nil {
			log.Fatalf("error on sending order %s", err)
		}
	}
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
