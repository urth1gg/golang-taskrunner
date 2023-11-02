package services

import (
	"caravagio-api-golang/internal/app/db"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type AuthService struct {
	authRepo db.AuthRepo
}

func NewAuthService(authRepo db.AuthRepo) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) ValidateAPIKey(ctx *gin.Context, apiKeyHeader string) (*db.APIKey, error) {
	if apiKeyHeader == "" {
		return &db.APIKey{}, errors.New("could not validate credentials")
	}

	apiKeyHeader = strings.Replace(apiKeyHeader, "Bearer ", "", 1)
	apiKey, err := s.authRepo.GetAPIKey(ctx, apiKeyHeader)
	if err != nil {
		return &db.APIKey{}, errors.New("could not validate credentials")
	}

	return apiKey, nil
}
