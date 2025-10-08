package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"
	"voucher-management-system/internal/dto"
	"voucher-management-system/internal/models"
	"voucher-management-system/internal/repository"

	"github.com/google/uuid"
)

type VoucherService interface {
	Create(dto.CreateVoucherDTO) (models.Voucher, error)
	FindAllPaginated(query dto.FindAllPaginatedVouchersDTO) ([]models.Voucher, int64, error)
	UploadCSV(file *multipart.FileHeader) (*CSVReport, error)
	GetAllVouchers() ([]models.Voucher, error)
	FindByID(uint) (models.Voucher, error)
	Update(uint, dto.UpdateVoucherDTO) (models.Voucher, error)
	Delete(uint) error
}

type CSVReport struct {
	TotalRows   int `json:"total_rows"`
	SuccessRows int `json:"success_rows"`
	FailedRows  int `json:"failed_rows"`
}

type voucherService struct {
	repo repository.VoucherRepository
}

func NewVoucherService(repo repository.VoucherRepository) VoucherService {
	return &voucherService{repo: repo}
}

func (s *voucherService) Create(data dto.CreateVoucherDTO) (models.Voucher, error) {
	layout := "2006-01-02" // format date
	expiryDate, errConvertDate := time.Parse(layout, data.ExpiryDate)
	if errConvertDate != nil {
		return models.Voucher{}, errors.New("invalid expiry_date format, must be YYYY-MM-DD")
	}

	voucher := models.Voucher{
		VoucherCode:     data.VoucherCode,
		DiscountPercent: data.DiscountPercent,
		ExpiryDate:      expiryDate,
	}
	err := s.repo.Create(&voucher)
	return voucher, err
}

func (s *voucherService) FindAllPaginated(query dto.FindAllPaginatedVouchersDTO) ([]models.Voucher, int64, error) {
	offset := (query.Page - 1) * query.Limit
	vouchers, err := s.repo.FindAll(query, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAll(query)
	if err != nil {
		return nil, 0, err
	}

	return vouchers, total, nil
}

func (s *voucherService) FindByID(id uint) (models.Voucher, error) {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return models.Voucher{}, err
	}
	return *voucher, nil
}

func (s *voucherService) Update(id uint, data dto.UpdateVoucherDTO) (models.Voucher, error) {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return models.Voucher{}, err
	}

	if data.DiscountPercent != nil {
		voucher.DiscountPercent = *data.DiscountPercent
	}

	if data.ExpiryDate != nil {
		layout := "2006-01-02"
		expiryDate, err := time.Parse(layout, *data.ExpiryDate)
		if err != nil {
			return models.Voucher{}, errors.New("invalid expiry_date format, must be YYYY-MM-DD")
		}
		voucher.ExpiryDate = expiryDate
	}

	if data.VoucherCode != nil {
		voucher.VoucherCode = *data.VoucherCode

	}

	if err := s.repo.Update(voucher); err != nil {
		return models.Voucher{}, err
	}

	return *voucher, nil
}

func (s *voucherService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *voucherService) UploadCSV(file *multipart.FileHeader) (*CSVReport, error) {
	// Simpan file sementara
	tempPath := fmt.Sprintf("/tmp/%s.csv", uuid.New().String())
	if err := saveTempFile(file, tempPath); err != nil {
		return nil, err
	}
	defer os.Remove(tempPath)

	f, err := os.Open(tempPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Validasi minimal 1 baris data
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must contain at least one data row")
	}

	// Hitung hasil
	report := &CSVReport{
		TotalRows: len(records) - 1, // tanpa header
	}

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		if len(row) < 3 {
			report.FailedRows++
			continue
		}

		voucher := models.Voucher{
			VoucherCode:     row[0],
			DiscountPercent: parseInt(row[1]),
			ExpiryDate:      parseDate(row[2]),
		}

		if err := s.repo.Create(&voucher); err != nil {
			report.FailedRows++
			continue
		}
		report.SuccessRows++
	}

	return report, nil
}

// Helper: simpan file sementara
func saveTempFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}

// Helper: parse integer
func parseInt(val string) int {
	var result int
	fmt.Sscanf(val, "%d", &result)
	return result
}

// Helper: parse date YYYY-MM-DD
func parseDate(dateStr string) time.Time {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dateStr)
	return t
}

func (s *voucherService) GetAllVouchers() ([]models.Voucher, error) {
	return s.repo.FindAllNoPagination()
}
