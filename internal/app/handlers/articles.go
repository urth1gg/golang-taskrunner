package handlers

import (
	"caravagio-api-golang/internal/app/models"
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ArticleHandler is the struct to hold necessary dependencies
type ArticleHandler struct {
	ArticleService   *services.ArticleService
	TaskQueueService *services.TaskQueueService
}

// HelloWorld responds with a "Hello, Articles!" message
func (h *ArticleHandler) HelloWorld(c *gin.Context) {
	value, errorfield := h.ArticleService.GetArticle(c, "74e13609-ee4e-4e74-812f-b1880cc732aa")

	c.JSON(200, gin.H{
		"message":    "Hello, Articles!",
		"value":      value,
		"errorfield": errorfield,
	})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	articleID := c.Param("articleID")

	value, errorfield := h.ArticleService.GetArticle(c, articleID)

	c.JSON(200, gin.H{
		"message":    "Hello, Articles!",
		"value":      value,
		"errorfield": errorfield,
	})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
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
		ID:          requestBody.Data.ID,
		Children:    requestBody.Data.Children,
		Tag:         requestBody.Data.Tag,
		Text:        requestBody.Data.Text,
		Keywords:    requestBody.Data.Keywords,
		PromptID:    requestBody.Data.PromptID,
		Expanded:    requestBody.Data.Expanded,
		IsCompleted: requestBody.Data.IsCompleted,
		Level:       requestBody.Data.Level,
		Length:      requestBody.Data.Length,
		MoreInfo:    requestBody.Data.MoreInfo,
		Response:    requestBody.Data.Response,
	}

	headingData.Data = append(headingData.Data, initialNode)

	article := models.Article{
		ArticleID:    articleID,
		HeadingData:  headingData,
		MainKeywords: requestBody.Data.Text,
	}

	h.ArticleService.UpdateArticle(c, &article)

	// Missplaced here but it's ok for now
	h.TaskQueueService.CreateTasksFromArticle(c, article)

	c.Data(http.StatusOK, "application/json", prettyJSON)
}

func (h *ArticleHandler) RegenerateHandler(c *gin.Context) {
	articleID := c.Param("articleID")

	article, err := h.ArticleService.GetArticle(c, articleID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.TaskQueueService.CreateTasksFromArticle(c, article)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully regenerated"})
}

// NewArticleHandler creates a new articles ArticleHandler
func NewArticleHandler(s *services.ArticleService, t *services.TaskQueueService) *ArticleHandler {
	return &ArticleHandler{ArticleService: s, TaskQueueService: t}
}
