package handlers

import (
	"TestTask/api-gateway/kafka"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"

	"net/http"
)

type Handlers struct {
	KafkaClient *kafka.KafkaClient
}

func (h *Handlers) CreateNewBalance(w http.ResponseWriter, r *http.Request) {
	var wallet kafka.Wallet

	if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wallet.ID = uuid.New().String()

	err := h.KafkaClient.ProduceNewWallet(wallet)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Send %v", wallet)
	status, err := h.KafkaClient.ConsumeStatus(wallet.ID)
	if err != nil {
		http.Error(w, "Error consuming status of the creation", http.StatusInternalServerError)
		return
	}

	switch status {
	case "error":
		http.Error(w, "Error during the execution of creation new wallet", http.StatusInternalServerError)
		return
	case "success":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successful creation")))
	}
}

func (h *Handlers) Invoice(w http.ResponseWriter, r *http.Request) {
	var transaction kafka.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction.ID = uuid.New().String()
	transaction.Type = "invoice"

	err := h.KafkaClient.ProduceTransaction(transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	status, err := h.KafkaClient.ConsumeStatus(transaction.ID)
	if err != nil {
		http.Error(w, "Error consuming status of the transaction", http.StatusInternalServerError)
		return
	}

	switch status {
	case "error":
		http.Error(w, "Error during the execution of the transaction", http.StatusInternalServerError)
		return
	case "success":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successful transaction")))
	}
}

func (h *Handlers) WithDraw(w http.ResponseWriter, r *http.Request) {
	var transaction kafka.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var err error

	transaction.ID = uuid.New().String()
	transaction.Type = "withdraw"

	err = h.KafkaClient.ProduceTransaction(transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	status, err := h.KafkaClient.ConsumeStatus(transaction.ID)
	if err != nil {
		http.Error(w, "Error consuming status of the transaction", http.StatusInternalServerError)
		return
	}

	switch status {
	case "error":
		http.Error(w, "Error during the execution of the transaction", http.StatusInternalServerError)
		return
	case "success":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successful transaction")))
	}
}

func (h *Handlers) ShowBalance(w http.ResponseWriter, r *http.Request) {
	//todo
}
