package kafka

import (
	"TestTask/transaction_system/internal/domain/db"
	"TestTask/transaction_system/internal/domain/models"
	"TestTask/transaction_system/pkg/client/postgres"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	Storage  db.BalanceStorage
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

	conDb, err := postgres.ConnectDB(context.Background())
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
		Storage: db.BalanceStorage{
			Context: context.Background(),
			Db:      conDb,
		},
	}, nil
}

func (kc *KafkaClient) StartWalletsConsumer(group *sync.WaitGroup) {
	defer group.Done()
	partitionConsumer, err := kc.consumer.ConsumePartition("wallets", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	log.Println("Wallets consumer starting")
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

			log.Printf("Received %v", wallet)

			if wallet.WalletNum == 0 {
				status := "created"

				if err := kc.Storage.CreateNewWallet(context.Background(), wallet); err != nil {
					status = "error"
					log.Printf("failed to create new wallet: %v", err)
				} else {
					status = "success"
				}

				statusJSON, err := json.Marshal(models.TransactionStatus{
					ID:     wallet.Id,
					Status: status,
				})
				if err != nil {
					log.Fatalf("failed to marshal status: %v", err)
				}
				kc.ProduceStatus("status", statusJSON)

			} else {
				status := "created"

				if err := kc.Storage.CreateNewCurrency(context.Background(), wallet); err != nil {
					status = "error"
					log.Printf("failed to create new wallet: %v", err)
				} else {
					status = "success"
				}

				statusJSON, err := json.Marshal(models.TransactionStatus{
					ID:     wallet.Id,
					Status: status,
				})
				if err != nil {
					log.Fatalf("failed to marshal status: %v", err)
				}
				kc.ProduceStatus("status", statusJSON)
			}
		}
	}
}

func (kc *KafkaClient) ProduceStatus(topic string, msg []byte) {
	_, _, err := kc.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		log.Fatalf("failed to send message to Kafka: %v", err)
	}
}

func (kc *KafkaClient) StartTransactionsConsumer(group *sync.WaitGroup) {
	defer group.Done()
	partitionConsumer, err := kc.consumer.ConsumePartition("transactions", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	log.Println("Transactions consumer starting")
	defer partitionConsumer.Close()

	for {
		select {
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}

			var transaction models.Transaction
			err := json.Unmarshal(msg.Value, &transaction)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			log.Printf("Received %v", transaction)

			if transaction.From == 0 {
				status := "created"

				if err := kc.Storage.Invoice(context.Background(), transaction); err != nil {
					status = "error"
					log.Printf("failed to invoice: %v", err)
				} else {
					status = "success"
				}

				statusJSON, err := json.Marshal(models.TransactionStatus{
					ID:     transaction.ID,
					Status: status,
				})
				if err != nil {
					log.Fatalf("failed to marshal status: %v", err)
				}
				kc.ProduceStatus("status", statusJSON)

			} else {
				status := "created"

				if err := kc.Storage.WithDraw(context.Background(), transaction); err != nil {
					status = "error"
					log.Printf("failed to withdraw: %v", err)
				} else {
					status = "success"
				}

				statusJSON, err := json.Marshal(models.TransactionStatus{
					ID:     transaction.ID,
					Status: status,
				})
				if err != nil {
					log.Fatalf("failed to marshal status: %v", err)
				}
				kc.ProduceStatus("status", statusJSON)
			}
		}
	}

}
