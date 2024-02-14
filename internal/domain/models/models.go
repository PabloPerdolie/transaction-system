package models

import "gorm.io/gorm"

type Transaction struct {
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Status   string  `json:"status"`
}

type Wallet struct {
	gorm.Model
	Id         int        `gorm:"primaryKey"`
	WalletNum  int        `gorm:"not null"`
	Currency   string     `gorm:"not null"`
	WalletData WalletData `gorm:"foreignKey:WalletID"`
}

type WalletData struct {
	gorm.Model
	WalletID int `gorm:"primaryKey"`
	//WalletNum     string  //`json:"wallet_num"`
	//Currency      string  `json:"currency"`
	ActualBalance float64 `gorm:"not null"` //`json:"actual"`
	FrozenBalance float64 `gorm:"not null"` //`json:"frozen"`
}
