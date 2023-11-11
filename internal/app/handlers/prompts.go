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
