package db

import (
	"TestTask/transaction_system/internal/domain/models"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

type BalanceStorage struct {
	Context context.Context
	Db      *gorm.DB
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

	result := bs.Db.Create(&wallet)
	if result.Error != nil {
		log.Printf("failed to create new wallet, %v", result.Error.Error())
		return result.Error
	}
	log.Printf("Created wallet: %v", wallet)
	return nil
}

func (bs *BalanceStorage) CreateNewCurrency(ctx context.Context, wallet models.Wallet) error {
	wallet = models.Wallet{
		WalletNum: wallet.WalletNum,
		Currency:  wallet.Currency,
		WalletData: models.WalletData{
			ActualBalance: 1000,
			FrozenBalance: 0,
		},
	}

	result := bs.Db.Create(&wallet)
	if result.Error != nil {
		log.Printf("failed to create new wallet, %v", result.Error.Error())
		return result.Error
	}
	log.Printf("Created wallet: %v", wallet)

	return nil
}

func (bs *BalanceStorage) Invoice(ctx context.Context, transaction models.Transaction) error {
	var wallet models.Wallet
	if err := bs.Db.Preload("WalletData").
		First(&wallet, "wallet_num = ? AND currency = ?", transaction.To, transaction.Currency).
		Error; err != nil {
		log.Printf("wallet doesn't exist, %v", err)
		return err
	}

	wallet.WalletData.ActualBalance += transaction.Amount

	if err := bs.Db.Save(&wallet).Error; err != nil {
		log.Printf("failed to invoice wallet, %v", err)
		return err
	}

	return nil
}

func (bs *BalanceStorage) WithDraw(ctx context.Context, transaction models.Transaction) error {
	var senderWallet models.Wallet
	if err := bs.Db.Preload("WalletData").
		First(&senderWallet, "wallet_num = ? AND currency = ?", transaction.From, transaction.Currency).
		Error; err != nil {
		log.Printf("sender's wallet doesn't exist, %v", err)
		return err
	}

	var receiverWallet models.Wallet
	if err := bs.Db.Preload("WalletData").
		First(&receiverWallet, "wallet_num = ? AND currency = ?", transaction.To, transaction.Currency).
		Error; err != nil {
		log.Printf("receiver's wallet doesn't exist, %v", err)
		return err
	}

	if senderWallet.WalletData.ActualBalance < transaction.Amount {
		return errors.New("Insufficient funds on the sender's account. ")
	}

	senderWallet.WalletData.ActualBalance -= transaction.Amount
	receiverWallet.WalletData.ActualBalance += transaction.Amount

	if err := bs.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&senderWallet).Error; err != nil {
			log.Printf("failed to save data , %v", err)
			return err
		}
		if err := tx.Save(&receiverWallet).Error; err != nil {
			log.Printf("failed to save data, %v", err)
			return err
		}
		return nil
	}); err != nil {
		log.Printf("failed to complete transaction, %v", err)
		return err
	}

	return nil
}

func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(10000000000000000)
	return randomNumber
}
