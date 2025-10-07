package main

import (
	"log"
	"voucher-management-system/config"
	"voucher-management-system/internal/models"
	"voucher-management-system/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	// Migrate models
	if err := config.DB.AutoMigrate(&models.Voucher{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")

}
