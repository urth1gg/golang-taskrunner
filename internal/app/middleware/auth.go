package middlewares

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthService interface {
	ValidateAPIKey(ctx *gin.Context, apiKey string) (db.APIKey, error)
}

type AuthMiddleware struct {
	service *services.AuthService
}

func NewAuthMiddleware(service *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKeyHeader := ctx.GetHeader("Authorization")

		if apiKey, err := m.service.ValidateAPIKey(ctx, apiKeyHeader); err == nil {
			ctx.Set("apiKey", apiKey)
			ctx.Next()
			return
		}

		cookie, err := ctx.Request.Cookie("api_key")

		if err == nil && cookie != nil {
			if apiKey, err := m.service.ValidateAPIKey(ctx, cookie.Value); err == nil {
				ctx.Set("Authorization", apiKey.Key)
				ctx.Next()
				return
			}

		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Could not validate credentials"})

	}
}
