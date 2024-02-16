package handlers

import (
	"TestTask/api-gateway/kafka"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) error {

	newClient, err := kafka.NewKafkaClient()
	if err != nil {
		return err
	}

	go newClient.StartConsumer()

	handlers := Handlers{KafkaClient: newClient}

	r.HandleFunc("/invoice", handlers.Invoice).Methods("POST")
	r.HandleFunc("/withdraw", handlers.WithDraw).Methods("POST")
	r.HandleFunc("/showbalance", handlers.ShowBalance).Methods("POST")
	r.HandleFunc("/create", handlers.CreateNewBalance).Methods("POST")

	return nil
}
