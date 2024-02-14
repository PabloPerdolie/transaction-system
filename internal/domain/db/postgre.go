package db

import (
	"TestTask/internal/domain/models"
	"context"
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
	wallet = models.Wallet{
		WalletNum: generateRandomNumber(),
		Currency:  wallet.Currency,
		WalletData: models.WalletData{
			ActualBalance: 1000,
			FrozenBalance: 0,
		},
	}

	result := bs.db.Create(&wallet)
	if result.Error != nil {
		log.Printf("failed to create new wallet, %v", result.Error.Error())
		return result.Error
	}

	return nil
}

func (bs *BalanceStorage) Invoice(ctx context.Context, transaction models.Transaction) error {
	var wallet models.Wallet
	if err := bs.db.Preload("WalletData").
		First(&wallet, "wallet_num = ? AND currency = ?", transaction.To, transaction.Currency).
		Error; err != nil {
		log.Printf("wallet doesn't exist, %v", err)
		return err
	}

	wallet.WalletData.ActualBalance += transaction.Amount

	if err := bs.db.Save(&wallet).Error; err != nil {
		log.Printf("failed to invoice wallet, %v", err)
		return err
	}

	return nil
}

func (bs *BalanceStorage) WithDraw(ctx context.Context, transaction models.Transaction) {

}

func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(10000000000000000)
	return randomNumber
}
