package handlers

import (
	"TestTask/transaction_system/internal/domain/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateNewBalance(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet

	if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//TODO доделать логику работы с кафкой

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("New wallet has been succecfully created")))
}

func Invoice(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

}

func WithDraw(w http.ResponseWriter, r *http.Request) {

}

func ShowBalance(w http.ResponseWriter, r *http.Request) {

}
