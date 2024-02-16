package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

type KafkaClient struct {
	producer             sarama.SyncProducer
	consumer             sarama.Consumer
	transactionStatusMap map[string]*TransactionStatus
	statusMutex          sync.Mutex
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
		producer:             producer,
		consumer:             consumer,
		transactionStatusMap: make(map[string]*TransactionStatus),
	}, nil
}

func (kc *KafkaClient) Close() {
	kc.producer.Close()
	kc.consumer.Close()
}

func (kc *KafkaClient) ProduceNewWallet(wallet Wallet) error {
	wallJSON, err := json.Marshal(wallet)
	if err != nil {
		log.Fatalf("failed to marshal wallet: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "wallets",
		Value: sarama.ByteEncoder(wallJSON),
	}

	_, _, err = kc.producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

func (kc *KafkaClient) ProduceTransaction(transaction Transaction) error {
	transJSON, err := json.Marshal(transaction)
	if err != nil {
		log.Fatalf("failed to marshal transaction: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "transactions",
		Value: sarama.ByteEncoder(transJSON),
	}

	_, _, err = kc.producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

func (kc *KafkaClient) ConsumeStatus(operationID string) (string, error) {
	kc.statusMutex.Lock()
	//statusChan := make(chan string)
	kc.transactionStatusMap[operationID] = &TransactionStatus{
		ID:        operationID,
		Status:    "",
		WaitGroup: sync.WaitGroup{},
	}
	kc.transactionStatusMap[operationID].WaitGroup.Add(1)
	kc.statusMutex.Unlock()

	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		defer kc.transactionStatusMap[operationID].WaitGroup.Done()

		select {
		case <-time.After(10 * time.Second):
		case <-doneChan:
		}
	}()
	log.Println("Wait......")
	kc.transactionStatusMap[operationID].WaitGroup.Wait()
	log.Println("Stop waiting!")
	kc.statusMutex.Lock()
	defer kc.statusMutex.Unlock()

	status := kc.transactionStatusMap[operationID].Status
	delete(kc.transactionStatusMap, operationID)

	if status != "" {
		return status, nil
	}

	return "", fmt.Errorf("waiting time is running out")
}

func (kc *KafkaClient) StartConsumer() {
	partitionConsumer, err := kc.consumer.ConsumePartition("status", 0, sarama.OffsetNewest)
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

			var statusMessage TransactionStatus
			err := json.Unmarshal(msg.Value, &statusMessage)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			kc.statusMutex.Lock()
			statusObj, exists := kc.transactionStatusMap[statusMessage.ID]
			if exists {
				statusObj.Status = statusMessage.Status
				statusObj.WaitGroup.Done()
			}
			kc.statusMutex.Unlock()
		}
	}
}
