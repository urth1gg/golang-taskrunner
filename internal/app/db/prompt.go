package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"fmt"
)

type PromptRepo interface {
	GetPrompt(ctx context.Context, promptID string) (*models.Prompt, error)
}

type DBPromptRepo struct {
	db *sql.DB
}

func NewDBPromptRepo(db *sql.DB) *DBPromptRepo {
	return &DBPromptRepo{db: db}
}

func (r *DBPromptRepo) GetPrompt(ctx context.Context, promptID string) (*models.Prompt, error) {
	var prompt models.Prompt
	err := r.db.QueryRowContext(ctx, "SELECT * FROM prompts WHERE prompt_id = ?", promptID).Scan(
		&prompt.PromptID,
		&prompt.UserID,
		&prompt.Name,
		&prompt.Description,
		&prompt.TextArea,
		&prompt.GPTModel,
		&prompt.Temperature,
		&prompt.MaxLength,
		&prompt.TopP,
		&prompt.FrequencyPenalty,
		&prompt.PresencePenalty,
		&prompt.CreatedAt,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &prompt, err
}
