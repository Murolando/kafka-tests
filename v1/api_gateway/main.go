package main

import (
	"api_gateway/handler"
	"api_gateway/server"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

var mu sync.Mutex

func main() {

	// Создание продюсера Kafka
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Создание консьюмера Kafka
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Подписка на партицию "pong" в Kafka
	partConsumer, err := consumer.ConsumePartition("oleg", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	// создание роутов для api
	h := handler.NewHandler(&producer)

	go func() {
		for {
			select {
			// Чтение сообщения из Kafka
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting goroutine")
					return
				}
				// когда мы прочитаем его, мы по id запихнем его в канал, чтобы горутина в hello поняла, что ответ вернулся
				// и вернул его сам домой
				fmt.Println("The message back from SecondMcr -> ApiGateway(main)",string(msg.Key))
				
				mu.Lock()
				responseID := string(msg.Key)
				h.Answer[responseID] = string(msg.Value)

				mu.Unlock()
			}
		}
	}()

	srv := new(server.Server)
	if err := srv.Run("8080", h.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
