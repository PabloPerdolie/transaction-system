package kafka

import (
	"TestTask/transaction_system/internal/domain/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"time"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewKafkaClient() (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 1000 * time.Millisecond
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		producer.Close()
		return nil, err
	}
	log.Println("Successfully connected to kafka and create consumer and producer")

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (kc *KafkaClient) StartConsumer() {
	partitionConsumer, err := kc.consumer.ConsumePartition("wallets", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	log.Println("Consumer starting")
	defer partitionConsumer.Close()

	for {
		select {
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}

			var wallet models.Wallet
			err := json.Unmarshal(msg.Value, &wallet)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			if wallet.WalletNum == 0 {
				log.Printf("Received %v", wallet)
				transJSON, err := json.Marshal(models.TransactionStatus{
					ID:     wallet.Id,
					Status: "success",
				})
				if err != nil {
					log.Fatalf("failed to marshal transaction: %v", err)
				}
				_, _, err = kc.producer.SendMessage(&sarama.ProducerMessage{
					Topic: "status",
					Value: sarama.ByteEncoder(transJSON),
				})
				if err != nil {
					log.Fatalf("failed to send message to Kafka: %v", err)
				}
			}
		}
	}
}
