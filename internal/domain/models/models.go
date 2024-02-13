package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Wallet   string  `json:"wallet"`
	Status   string  `json:"status"`
}

type Balance struct {
	gorm.Model
	Wallet   string  `json:"wallet"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Frozen   float64 `json:"frozen"`
}
