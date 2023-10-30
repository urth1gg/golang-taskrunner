package services

import (
	"errors"
	"strings"
	"caravagio-api-golang/internal/app/db"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	authRepo db.AuthRepo
}

func NewAuthService(authRepo db.AuthRepo) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) ValidateAPIKey(ctx *gin.Context, apiKeyHeader string) (string, error) {
	if apiKeyHeader == "" {
		return "", errors.New("could not validate credentials")
	}

	apiKeyHeader = strings.Replace(apiKeyHeader, "Bearer ", "", 1)
	apiKey, err := s.authRepo.GetAPIKey(ctx, apiKeyHeader)
	if err != nil {
		return "", errors.New("could not validate credentials")
	}

	return apiKey.Key, nil
}