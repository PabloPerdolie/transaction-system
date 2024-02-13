package postgres

import (
	"TestTask/internal/domain/models"
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func ConnectDB(ctx context.Context) (*gorm.DB, error) {

	dsn := os.Getenv("PSQL_AUTH")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	} else {
		log.Println("Successfully connected to database")
	}

	if err = db.AutoMigrate(&models.Wallet{}); err != nil {
		log.Fatalf("Failed to AutoMigrate: %v", err)
		return nil, err
	}

	return db, nil
}
