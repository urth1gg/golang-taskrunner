package handlers

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SettingsHandler struct {
	settingsService       *services.SettingsService
	defaultPromptsService *services.DefaultPromptsService
}

func NewSettingsHandler(settingsService *services.SettingsService, defaultPromptsService *services.DefaultPromptsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService, defaultPromptsService: defaultPromptsService}
}

type DefaultPrompt struct {
	UserID                 string `sql:"user_id"`
	PromptID               string `sql:"prompt_id"`
	HeadingNameAndPosition string `sql:"heading_name_and_position"`
}

func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID := c.Param("userID")

	defaultPrompts, err := h.defaultPromptsService.GetAllDefaultPrompts(c, userID)

	if err != nil {
		fmt.Println(err)
		defaultPrompts = []db.DefaultPrompt{}
	}

	c.JSON(200, gin.H{
		"defaultPrompts": defaultPrompts,
	})
}

func (h *SettingsHandler) UpdateDefaultPrompts(c *gin.Context) {
	var requestBody map[string][]map[string]string // Adjusted type to match the JSON structure

	userID := c.Param("userID")
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prettyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Accessing prompts
	prompts := requestBody["prompts"]

	defaultPrompts := []db.DefaultPrompt{}

	// Iterating over the array of maps
	for _, prompt := range prompts {
		for key, value := range prompt {

			defaultPrompt := db.DefaultPrompt{
				UserID:                 userID,
				PromptID:               value,
				HeadingNameAndPosition: key,
			}

			defaultPrompts = append(defaultPrompts, defaultPrompt)
		}
	}

	for _, defaultPrompt := range defaultPrompts {
		_, err := h.defaultPromptsService.UpdateDefaultPrompt(c, &defaultPrompt)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{
		"message":    "Settings updated successfully",
		"prettyJSON": string(prettyJSON),
	})
}

