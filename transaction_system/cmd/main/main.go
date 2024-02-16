package main

import (
	"TestTask/transaction_system/internal/domain/kafka"
)

func main() {
	client, err := kafka.NewKafkaClient()
	if err != nil {
		panic(err)
	}

	client.StartConsumer()
}
