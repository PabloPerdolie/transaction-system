package kafka

import "sync"

type Wallet struct {
	ID         string
	WalletNum  int        `json:"wallet_num"`
	Currency   string     `json:"cur"`
	WalletData WalletData `json:"wallet_data"`
}

type WalletData struct {
	ActualBalance float64 `json:""` //todo
	FrozenBalance float64 `json:""`
}

type TransactionStatus struct {
	ID        string
	Status    string
	WaitGroup sync.WaitGroup
}

type Transaction struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Status   string  `json:"status"`
}
