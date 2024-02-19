package kafka

import "sync"

type Wallet struct {
	Id        string
	WalletNum int    `json:"wallet_num"`
	Currency  string `json:"cur"`
	//WalletData WalletData `json:"wallet_data"`
}

//type WalletData struct {
//	ActualBalance float64 `json:""` //todo
//	FrozenBalance float64 `json:""`
//}

type TransactionStatus struct {
	ID        string
	Status    string
	WaitGroup sync.WaitGroup
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
