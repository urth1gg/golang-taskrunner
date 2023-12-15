package handlers

import (
	"caravagio-api-golang/internal/app/models"
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"fmt"
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

	if requestBody.FinishSentence {
		h.TaskQueueService.CreateFinishSentenceTasksFromArticle(c, article)
	}

	if requestBody.FixGrammar {
		fmt.Printf("%v", article)
		h.TaskQueueService.CreateFixGrammarTasksFromArticle(c, article)
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
/*
INSERT INTO prompts (
    prompt_id,
    user_id,
    name,
    description,
    text_area,
    gpt_model,
    temperature,
    max_length,
    top_p,
    frequency_penalty,
    presence_penalty,
    created_at,
    level_required_to_access
)
VALUES (
    'bd56f391-9ae6-11ee-8fe2-00155d509f69',
    'eb360df1-d5a6-467b-8bda-7ba02e114e4c',
    'Fix grammar name',
    'Fix grammar desc',
    'Fix any grammatical errors you find in the text below: \n\n {text}',
    'gpt-4-1106-preview',
    0.8,
    1000,
    1,
    0.3,
    0.3,
    '2023-12-15 02:10:30',
    '1'
);

INSERT INTO prompts (
    prompt_id,
    user_id,
    name,
    description,
    text_area,
    gpt_model,
    temperature,
    max_length,
    top_p,
    frequency_penalty,
    presence_penalty,
    created_at,
    level_required_to_access
)
VALUES (
    '6727d92b-9ae6-11ee-8fe2-00155d509f69',
    'eb360df1-d5a6-467b-8bda-7ba02e114e4c',
    'Finish sentence name',
    'Finish sentence desc',
    'Finish the sentence and paragraph that you’ve started but don’t write more: \n\n {text}',
    'gpt-4-1106-preview',
    0.8,
    1000,
    1,
    0.3,
    0.3,
    '2023-12-15 02:08:05',
    '1'
);
*/
