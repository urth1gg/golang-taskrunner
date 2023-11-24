package handlers

import (
	"caravagio-api-golang/internal/app/services"
	"context"
)

type PromptsHandler struct {
	promptService *services.PromptService
}

func (h *PromptsHandler) GetPrompt(promptID string) (string, error) {

	ctx := context.Background()
	prompt, err := h.promptService.GetPrompt(ctx, promptID)

	if err != nil {
		return "", err
	}

	return prompt.TextArea.String, nil
}

func (h *PromptsHandler) GetAllAvailablePrompts(levelRequiredToAccess string) ([]string, error) {

	ctx := context.Background()
	prompts, err := h.promptService.GetAllAvailablePrompts(ctx, levelRequiredToAccess)

	if err != nil {
		return []string{}, err
	}

	var promptsText []string

	for _, prompt := range prompts {
		promptsText = append(promptsText, prompt.TextArea.String)
	}

	return promptsText, nil
}

func NewPromptsHandler(promptService *services.PromptService) *PromptsHandler {
	return &PromptsHandler{promptService: promptService}
}
