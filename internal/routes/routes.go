package routes

import (
	"voucher-management-system/internal/middleware"
	"voucher-management-system/internal/repository"
	"voucher-management-system/internal/services"

	"voucher-management-system/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Auth setup
	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService)

	voucherRepo := repository.NewVoucherRepository()
	voucherService := services.NewVoucherService(voucherRepo)
	voucherHandler := handlers.NewVoucherHandler(voucherService)

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/test", voucherHandler.FindAllPaginatedVouchers)
		voucher := api.Group("/vouchers", middleware.AuthMiddleware())
		{
			voucher.POST("", voucherHandler.CreateVoucher)
			voucher.GET("", voucherHandler.FindAllPaginatedVouchers)
			voucher.GET(":id", voucherHandler.GetVoucherByID)
			voucher.PUT(":id", voucherHandler.UpdateVoucher)
			voucher.DELETE(":id", voucherHandler.DeleteVoucher)
			voucher.POST("upload-csv", voucherHandler.UploadCSV)
			voucher.POST("export", voucherHandler.ExportVouchers)
		}
	}
}
