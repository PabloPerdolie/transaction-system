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
	WalletID      string  `gorm:"uniqueIndex;autoIncrement"`
	WalletNum     string  //`json:"wallet_num"`
	Currency      string  `json:"currency"`
	ActualBalance float64 `json:"amount"`
	FrozenBalance float64 `json:"frozen"`
}
