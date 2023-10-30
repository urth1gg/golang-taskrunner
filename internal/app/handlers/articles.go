package articleshandler

import (
	"github.com/gin-gonic/gin"
	"caravagio-api-golang/internal/app/services"
	"net/http"
	"encoding/json"
	// "fmt"
	"caravagio-api-golang/internal/app/models"
)

// Handler is the struct to hold necessary dependencies
type Handler struct {
	ArticleService *services.ArticleService
}

// HelloWorld responds with a "Hello, Articles!" message
func (h *Handler) HelloWorld(c *gin.Context) {
	value, errorfield := h.ArticleService.GetArticle(c, "74e13609-ee4e-4e74-812f-b1880cc732aa")

	c.JSON(200, gin.H{
		"message": "Hello, Articles!",
		"value": value,
		"errorfield": errorfield,
	})
}

func (h *Handler) GetArticle(c *gin.Context) {
	articleID := c.Param("articleID")

	value, errorfield := h.ArticleService.GetArticle(c, articleID)

	c.JSON(200, gin.H{
		"message": "Hello, Articles!",
		"value": value,
		"errorfield": errorfield,
	})
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	var requestBody models.ArticleBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prettyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", prettyJSON)
}

// NewHandler creates a new articles handler
func NewHandler(s *services.ArticleService) *Handler {
	return &Handler{ArticleService: s}
}

