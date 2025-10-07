package repository

import (
	"voucher-management-system/config"
	"voucher-management-system/internal/dto"
	"voucher-management-system/internal/models"
)

type VoucherRepository interface {
	Create(voucher *models.Voucher) error
	FindAll(query dto.FindAllPaginatedVouchersDTO, offset int) ([]models.Voucher, error)
	CountAll(query dto.FindAllPaginatedVouchersDTO) (int64, error)
	FindByID(id uint) (*models.Voucher, error)
	Update(voucher *models.Voucher) error
	Delete(id uint) error
	FindAllNoPagination() ([]models.Voucher, error)
}

type voucherRepository struct{}

func NewVoucherRepository() VoucherRepository {
	return &voucherRepository{}
}

func (r *voucherRepository) Create(voucher *models.Voucher) error {
	return config.DB.Create(voucher).Error
}

func (r *voucherRepository) FindAll(query dto.FindAllPaginatedVouchersDTO, offset int) ([]models.Voucher, error) {
	var vouchers []models.Voucher

	db := config.DB.Model(&models.Voucher{})
	if query.Search != "" {
		db = db.Where("voucher_code ILIKE ?", "%"+query.Search+"%")
	}

	err := db.Limit(query.Limit).Offset(offset).Order("id DESC").Find(&vouchers).Error
	return vouchers, err
}
func (r *voucherRepository) FindAllNoPagination() ([]models.Voucher, error) {
	var vouchers []models.Voucher
	err := config.DB.Order("id asc").Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepository) CountAll(query dto.FindAllPaginatedVouchersDTO) (int64, error) {
	var total int64

	db := config.DB.Model(&models.Voucher{})
	if query.Search != "" {
		db = db.Where("voucher_code ILIKE ?", "%"+query.Search+"%")
	}

	err := db.Count(&total).Error
	return total, err
}

func (r *voucherRepository) FindByID(id uint) (*models.Voucher, error) {
	var voucher models.Voucher
	err := config.DB.First(&voucher, id).Error
	return &voucher, err
}

func (r *voucherRepository) Update(voucher *models.Voucher) error {
	return config.DB.Save(voucher).Error
}

func (r *voucherRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Voucher{}, id).Error
}
