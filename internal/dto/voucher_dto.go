package dto

import "time"

type CreateVoucherDTO struct {
	VoucherCode     string `json:"voucher_code" binding:"required"`
	DiscountPercent int    `json:"discount_percent" binding:"required,gte=0,lte=100"`
	ExpiryDate      string `json:"expiry_date" binding:"required"`
}

type FindAllPaginatedVouchersDTO struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
}

type UpdateVoucherDTO struct {
	DiscountPercent *int    `json:"discount_percent" binding:"omitempty,gte=0,lte=100"`
	ExpiryDate      *string `json:"expiry_date" binding:"omitempty"`
}

type VoucherResponse struct {
	ID              uint      `json:"id"`
	No              uint      `json:"no"`
	VoucherCode     string    `json:"voucher_code"`
	DiscountPercent int       `json:"discount_percent"`
	ExpiryDate      time.Time `json:"expiry_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
