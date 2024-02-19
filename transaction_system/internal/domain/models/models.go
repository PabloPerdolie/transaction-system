package models

import (
	"gorm.io/gorm"
)

type TransactionStatus struct {
	ID     string
	Status string
	//WaitGroup *sync.WaitGroup
}

type Transaction struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Currency string  `json:"cur"`
	Amount   float64 `json:"amount"`
	From     int     `json:"from"`
	To       int     `json:"to"`
	Status   string  `json:"status"`
}

type Wallet struct {
	gorm.Model
	Id         string     `gorm:"-"`
	WalletNum  int        `json:"wallet_num" ,gorm:"not null"`
	Currency   string     `json:"cur" ,gorm:"not null"`
	WalletId   int        `gorm:"not null"`
	WalletData WalletData //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type WalletData struct {
	gorm.Model
	ID            int     `gorm:"primaryKey"`
	ActualBalance float64 `gorm:"not null"` //`json:"actual"`
	FrozenBalance float64 `gorm:"not null"` //`json:"frozen"`
	WalletId      int     `gorm:"autoIncrement;autoIncrement:1"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	walletData := WalletData{
		ActualBalance: 1000,
		FrozenBalance: 0,
	}
	if err := tx.Create(&walletData).Error; err != nil {
		return err
	}

	w.WalletId = walletData.WalletId

	return nil
}
