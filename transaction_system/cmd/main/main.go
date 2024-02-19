package main

import (
	"TestTask/transaction_system/internal/domain/kafka"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		panic(err)
	}
	client, err := kafka.NewKafkaClient()
	if err != nil {
		panic(err)
	}
	var WaitGroup sync.WaitGroup
	WaitGroup.Add(2)
	go client.StartWalletsConsumer(&WaitGroup)
	go client.StartTransactionsConsumer(&WaitGroup)
	WaitGroup.Wait()
}
