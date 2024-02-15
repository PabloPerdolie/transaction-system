package handlers

import (
	"TestTask/transaction_system/internal/adapters/api/balance"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/invoice", balance.Invoice).Methods("POST")
	r.HandleFunc("/withdraw", balance.WithDraw).Methods("POST")
	r.HandleFunc("/showbalance", balance.ShowBalance).Methods("POST")
	r.HandleFunc("/create", balance.CreateNewBalance).Methods("POST")
}
