package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthService interface {
	ValidateAPIKey(ctx *gin.Context, apiKey string) (string, error)
}

type AuthMiddleware struct {
	service AuthService
}

func NewAuthMiddleware(service AuthService) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKeyHeader := ctx.GetHeader("Authorization")
		
		if apiKey, err := m.service.ValidateAPIKey(ctx, apiKeyHeader); err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Could not validate credentials"})
			return
		} else {
			ctx.Set("apiKey", apiKey)
			ctx.Next()
		}
	}
}