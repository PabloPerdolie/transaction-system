package models

import (
	"gorm.io/gorm"
	"time"
)

type TransactionStatus struct {
	ID     string
	Status string
	//WaitGroup *sync.WaitGroup
}

type Transaction struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	Amount    float64   `json:"amount"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"` //todo
}

type Wallet struct {
	gorm.Model
	Id         string     `gorm:"-"`
	WalletNum  int        `json:"wallet_num" ,gorm:"not null"`
	Currency   string     `json:"cur" ,gorm:"not null"`
	WalletData WalletData `gorm:"foreignKey:WalletID"`
}

type WalletData struct {
	gorm.Model
	WalletID      int     `gorm:"primaryKey"`
	ActualBalance float64 `gorm:"not null"` //`json:"actual"`
	FrozenBalance float64 `gorm:"not null"` //`json:"frozen"`
}
