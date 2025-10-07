package main

import (
	"log"
	"voucher-management-system/config"
	"voucher-management-system/internal/models"
)

func main() {
	log.Println("Starting database migration...")

	// DATABASE CONNECTION
	config.ConnectDatabase()

	// RUN AUTO MIGRATE FOR ALL MODELS
	err := config.DB.AutoMigrate(
		&models.Voucher{},
		&models.Authentication{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
