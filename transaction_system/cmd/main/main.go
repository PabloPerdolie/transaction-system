package main

import (
	"TestTask/transaction_system/internal/domain/kafka"
	"github.com/joho/godotenv"
	"log"
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

	client.StartConsumer()
}
