package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"fmt"
)

type PromptRepo interface {
	GetPrompt(ctx context.Context, promptID string) (*models.Prompt, error)
	GetAllAvailablePrompts(ctx context.Context, levelRequiredToAccess string) ([]models.Prompt, error)
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
		&prompt.LevelRequiredToAccess,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &prompt, err
}

func (r *DBPromptRepo) GetAllAvailablePrompts(ctx context.Context, levelRequiredToAccess string) ([]models.Prompt, error) {
	var prompts []models.Prompt
	query := "SELECT * FROM prompts WHERE level_required_to_access <= ?"
	rows, err := r.db.QueryContext(ctx, query, levelRequiredToAccess)

	if err != nil {
		fmt.Println(err)
		return prompts, err
	}

	for rows.Next() {
		var prompt models.Prompt
		err = rows.Scan(
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
			&prompt.LevelRequiredToAccess,
		)

		if err != nil {
			fmt.Println(err)
			return prompts, err
		}

		// Don't return prompts Text Area to the user if under certain level

		if levelRequiredToAccess == "1" {
			prompt.TextArea = sql.NullString{}
		}

		prompts = append(prompts, prompt)
	}

	return prompts, nil
}
