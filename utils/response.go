package utils

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// StandardResponse adalah format respons global
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type MetaResponse struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	TotalData  int64 `json:"total_data,omitempty"`
	TotalPages int64 `json:"total_pages,omitempty"`
}

// SuccessResponse untuk respons sukses
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	response := gin.H{
		"success": true,
		"message": message,
		"data":    data,
	}

	c.JSON(code, response)
}
func SuccessResponsePagination(c *gin.Context, code int, message string, data interface{}, meta *MetaResponse) {
	response := gin.H{
		"success": true,
		"message": message,
		"data":    data,
	}

	// only show meta if not nil
	if meta != nil {
		response["meta"] = meta
	}

	c.JSON(code, response)
}

func ErrorResponse(c *gin.Context, code int, message string, err error) {
	errorMessage := err.Error()
	duplicateField := ""

	// Detect error duplicate key PostgreSQL
	if strings.Contains(errorMessage, "duplicate key value violates unique constraint") {
		// Get nama constraint dari error
		re := regexp.MustCompile(`unique constraint "([^"]+)"`)
		match := re.FindStringSubmatch(errorMessage)
		if len(match) > 1 {
			constraint := match[1]
			// Usually constraint formatnya "idx_<table>_<field>"
			parts := strings.Split(constraint, "_")
			if len(parts) >= 3 {
				duplicateField = parts[len(parts)-1]
			}
		}
		errorMessage = "Data already exists"
		code = http.StatusConflict
	}

	// Detect error not found
	if strings.Contains(errorMessage, "record not found") {
		errorMessage = "Data not found"
		code = http.StatusNotFound
	}

	// Base response
	response := gin.H{
		"success": false,
		"message": message,
		"error":   errorMessage,
	}

	// If ada duplicate field, add key `field`
	if duplicateField != "" {
		response["field"] = duplicateField
	}

	c.JSON(code, response)
}

// ValidationErrorResponse untuk validasi khusus
func ValidationErrorResponse(ctx *gin.Context, err error) {
	formatted := FormatValidationError(err)
	ctx.JSON(http.StatusBadRequest, StandardResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  formatted,
	})
}
