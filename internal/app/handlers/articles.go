package articleshandler

import (
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
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

	articleID := c.Param("articleID")
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prettyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	headingData := models.HeadingData{
		Data: []models.Node{},
	}

	initialNode := models.Node{
		ID: requestBody.Data.ID,
		Children: requestBody.Data.Children,
		Tag: requestBody.Data.Tag,
		Text: requestBody.Data.Text,
		Keywords: requestBody.Data.Keywords,
		PromptID: requestBody.Data.PromptID,
		Expanded: requestBody.Data.Expanded,
		IsCompleted: requestBody.Data.IsCompleted,
		Level: requestBody.Data.Level,
		Length: requestBody.Data.Length,
		MoreInfo: requestBody.Data.MoreInfo,
		Response: requestBody.Data.Response,
	}

	headingData.Data = append(headingData.Data, initialNode)

	article := models.Article{
		ArticleID: articleID,
		HeadingData: headingData,
	}

	log.Println("ArticleID")
	log.Println(article.ArticleID)
	h.ArticleService.UpdateArticle(c, &article)

	c.Data(http.StatusOK, "application/json", prettyJSON)
}

// NewHandler creates a new articles handler
func NewHandler(s *services.ArticleService) *Handler {
	return &Handler{ArticleService: s}
}

