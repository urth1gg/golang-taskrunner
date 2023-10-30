package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
)

type PromptService struct {
	db db.PromptRepo
}

func NewPromptService(db db.PromptRepo) *PromptService {
	return &PromptService{db: db}
}

func (s *PromptService) GetPrompt(ctx context.Context, promptID string) (models.Prompt, error) {
	// get first row from prompts table based on promptID

	prompt, err := s.db.GetPrompt(ctx, promptID)

	if err != nil {
		fmt.Println(err)
		return models.Prompt{}, err
	}

	return *prompt, nil
}
