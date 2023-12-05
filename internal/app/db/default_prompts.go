package db

import (
	"context"
	"database/sql"
)

type DefaultPromptsRepo interface {
	GetDefaultPrompt(ctx context.Context, promptID string, userID string) (*DefaultPrompt, error)
	GetAllDefaultPrompts(ctx context.Context, userID string) ([]DefaultPrompt, error)
	UpdateDefaultPrompt(ctx context.Context, prompt *DefaultPrompt) (int, error)
	CreateDefaultPrompt(ctx context.Context, prompt *DefaultPrompt) (*DefaultPrompt, error)
}

type DBDefaultPromptsRepo struct {
	db *sql.DB
}

type DefaultPrompt struct {
	UserID                 string `sql:"user_id"`
	PromptID               string `sql:"prompt_id"`
	HeadingNameAndPosition string `sql:"heading_name_and_position"`
}

func NewDBDefaultPromptsRepo(db *sql.DB) *DBDefaultPromptsRepo {
	return &DBDefaultPromptsRepo{db: db}
}

func (r *DBDefaultPromptsRepo) GetDefaultPrompt(ctx context.Context, promptID string, userID string) (*DefaultPrompt, error) {
	var prompt DefaultPrompt
	query := "SELECT user_id, prompt_id, heading_name_and_position FROM default_prompts WHERE prompt_id = ? AND user_id = ?"
	err := r.db.QueryRowContext(ctx, query, promptID, userID).Scan(
		&prompt.UserID,
		&prompt.PromptID,
		&prompt.HeadingNameAndPosition,
	)

	if err != nil {
		return nil, err
	}

	return &prompt, nil
}

func (r *DBDefaultPromptsRepo) GetAllDefaultPrompts(ctx context.Context, userID string) ([]DefaultPrompt, error) {
	var prompts []DefaultPrompt
	query := "SELECT user_id, prompt_id, heading_name_and_position FROM default_prompts WHERE user_id = ?"
	rows, err := r.db.QueryContext(ctx, query, userID)

	if err != nil {
		return prompts, err
	}

	for rows.Next() {
		var prompt DefaultPrompt
		err = rows.Scan(
			&prompt.UserID,
			&prompt.PromptID,
			&prompt.HeadingNameAndPosition,
		)
		if err != nil {
			return prompts, err
		}
		prompts = append(prompts, prompt)
	}

	return prompts, nil
}

func (r *DBDefaultPromptsRepo) UpdateDefaultPrompt(ctx context.Context, prompt *DefaultPrompt) (int, error) {
	query := `INSERT INTO default_prompts (user_id, heading_name_and_position, prompt_id)
              VALUES (?, ?, ?)
              ON DUPLICATE KEY UPDATE 
                  prompt_id = VALUES(prompt_id)`

	result, err := r.db.ExecContext(ctx, query, prompt.UserID, prompt.HeadingNameAndPosition, prompt.PromptID)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, nil
	}

	return int(rowsAffected), nil
}

func (r *DBDefaultPromptsRepo) CreateDefaultPrompt(ctx context.Context, prompt *DefaultPrompt) (*DefaultPrompt, error) {
	query := "INSERT INTO default_prompts (user_id, prompt_id, heading_name_and_position) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, prompt.UserID, prompt.PromptID, prompt.HeadingNameAndPosition)

	if err != nil {
		return nil, err
	}

	return prompt, nil
}
