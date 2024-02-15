package kafka

import (
	"TestTask/transaction_system/internal/domain/models"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

var (
	transactionStatusMap = make(map[string]*models.TransactionStatus)
	statusMutex          sync.Mutex
	producer             sarama.SyncProducer
	consumer             sarama.Consumer
)

func init() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 1000 * time.Millisecond

	brokers := []string{"localhost:9092"}

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}

	consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
}

func produceTransaction(transaction models.Transaction) error {
	return nil
}

func consumeTransactionStatus(transactionID string) (string, error) {
	return "", nil
}
