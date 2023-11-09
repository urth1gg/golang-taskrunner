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
		Keywords:     requestBody.Data.Keywords,
		MoreInfo:     requestBody.Data.MoreInfo,
	}

	fieldsToUpdate := []string{"heading_data", "main_keywords", "keywords"}
	h.ArticleService.UpdateArticleGeneric(c, &article, fieldsToUpdate)

	if requestBody.Regenerate {
		tasksInProgress, err := h.TaskQueueService.GetAllInProgressTasksByArticleId(c, &article)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.TaskQueueService.CancelResponseStreamForTasks(c, &tasksInProgress)
		h.TaskQueueService.DeleteTasksByArticleId(c, &article)
		h.TaskQueueService.CreateTasksFromArticle(c, article)
	}

	if requestBody.Continue {
		h.TaskQueueService.CreateContinueTasksFromArticle(c, article)
	}

	if requestBody.MetaDescription.ID != "" {
		_, err := h.TaskQueueService.CreateMetaDescriptionTask(c, &article, &requestBody.MetaDescription)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

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

func (h *ArticleHandler) DeleteTasks(c *gin.Context) {

	h.TaskQueueService.DeleteTasks(c)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted tasks"})
}

// NewArticleHandler creates a new articles ArticleHandler
func NewArticleHandler(s *services.ArticleService, t *services.TaskQueueService) *ArticleHandler {
	return &ArticleHandler{ArticleService: s, TaskQueueService: t}
}

//TODO: Delete this
//ALTER TABLE articles ADD COLUMN meta_description TEXT DEFAULT '';
