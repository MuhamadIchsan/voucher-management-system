package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"voucher-management-system/internal/dto"
	"voucher-management-system/internal/services"
	"voucher-management-system/utils"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	service services.VoucherService
}

func NewVoucherHandler(service services.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: service}
}

// CreateVoucher godoc
func (h *VoucherHandler) CreateVoucher(c *gin.Context) {
	var body dto.CreateVoucherDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	voucher, err := h.service.Create(body)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create voucher", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Voucher created successfully", voucher)
}

// FindAllPaginatedVouchers godoc
func (h *VoucherHandler) FindAllPaginatedVouchers(c *gin.Context) {
	var query dto.FindAllPaginatedVouchersDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	vouchers, total, err := h.service.FindAllPaginated(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch vouchers", err)
		return
	}

	totalPages := utils.CountPages(total, int64(query.Limit))
	meta := &utils.MetaResponse{
		Page:       query.Page,
		Limit:      query.Limit,
		TotalData:  total,
		TotalPages: totalPages,
	}

	utils.SuccessResponsePagination(c, http.StatusOK, "Vouchers retrieved successfully", vouchers, meta)
}

// GetVoucherByID godoc
func (h *VoucherHandler) GetVoucherByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID", err)
		return
	}

	voucher, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Voucher not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Voucher retrieved successfully", voucher)
}

// UpdateVoucher godoc
func (h *VoucherHandler) UpdateVoucher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID", err)
		return
	}

	var body dto.UpdateVoucherDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	voucher, err := h.service.Update(uint(id), body)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update voucher", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Voucher updated successfully", voucher)
}

// DeleteVoucher godoc
func (h *VoucherHandler) DeleteVoucher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID", err)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete voucher", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Voucher deleted successfully", nil)
}

// UploadCSV godoc
func (h *VoucherHandler) UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CSV file is required", err)
		return
	}

	report, err := h.service.UploadCSV(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process CSV file", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "CSV processed successfully", report)
}

func (h *VoucherHandler) ExportVouchers(c *gin.Context) {
	vouchers, err := h.service.GetAllVouchers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch vouchers", err)
		return
	}

	// Set headers for CSV download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="vouchers.csv"`)
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Voucher Code", "Discount Percent", "Expiry Date", "Created At"})

	// Write data
	for _, v := range vouchers {
		writer.Write([]string{
			strconv.FormatUint(uint64(v.ID), 10),
			v.VoucherCode,
			strconv.Itoa(v.DiscountPercent),
			v.ExpiryDate.Format("2006-01-02"),
			v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}
