package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Dummy check — in real app you'd query DB
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	// Generate random token
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenBytes), nil
}
