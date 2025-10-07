package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Voucher struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	No              uint           `gorm:"uniqueIndex;autoIncrement" json:"no"`
	VoucherCode     string         `gorm:"uniqueIndex;not null" json:"voucher_code"`
	DiscountPercent int            `gorm:"not null;default:0" json:"discount_percent"`
	ExpiryDate      time.Time      `gorm:"type:date;not null" json:"expiry_date"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Voucher) BeforeSave(tx *gorm.DB) (err error) {
	if v.DiscountPercent < 0 || v.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	return nil
}
