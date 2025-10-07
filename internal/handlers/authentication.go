package handlers

import (
	"net/http"
	"voucher-management-system/internal/services"
	"voucher-management-system/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SuccessResponse(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		utils.SuccessResponse(c, http.StatusUnauthorized, "Invalid request", gin.H{"error": err.Error()})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login Success", token)

}
