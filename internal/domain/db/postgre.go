package db

import (
	"TestTask/internal/domain/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

type BalanceStorage struct {
	context context.Context
	db      *gorm.DB
}

//func NewBalanceStorage() BalanceStorage {
//	return BalanceStorage{
//
//	}
//}

func (bs *BalanceStorage) CreateNewWallet(ctx context.Context, wallet models.Wallet) error {
	wallet.WalletNum = generateRandomNumber()
	wallet.ActualBalance = 0
	wallet.FrozenBalance = 0

	result := bs.db.Create(&wallet)
	if result.Error != nil {
		log.Printf("failed to create new wallet, %v", result.Error.Error())
		return result.Error
	}

	return nil
}

func (bs *BalanceStorage) Invoice(ctx context.Context, transaction models.Transaction) error {

	return nil
}

func generateRandomNumber() string {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(10000000000000000)
	return fmt.Sprintf("%016d", randomNumber)
}
