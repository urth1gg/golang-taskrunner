package services

import (
	"caravagio-api-golang/internal/app/db"
	"context"
	"fmt"
)

type DefaultPromptsService struct {
	db db.DefaultPromptsRepo
}

func NewDefaultPromptsService(db db.DefaultPromptsRepo) *DefaultPromptsService {
	return &DefaultPromptsService{db: db}
}

func (s *DefaultPromptsService) GetDefaultPrompt(ctx context.Context, promptID string, userID string) (*db.DefaultPrompt, error) {
	// get first row from prompts table based on promptID

	prompt, err := s.db.GetDefaultPrompt(ctx, promptID, userID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prompt, nil
}

func (s *DefaultPromptsService) GetAllDefaultPrompts(ctx context.Context, userID string) ([]db.DefaultPrompt, error) {
	// get first row from prompts table based on promptID

	prompts, err := s.db.GetAllDefaultPrompts(ctx, userID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prompts, nil
}

func (s *DefaultPromptsService) UpdateDefaultPrompt(ctx context.Context, prompt *db.DefaultPrompt) (int, error) {

	_, err := s.db.UpdateDefaultPrompt(ctx, prompt)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return 1, nil
}
