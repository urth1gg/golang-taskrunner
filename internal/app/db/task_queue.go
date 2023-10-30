package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"fmt"
)

type TaskQueueRepo interface {
	GetTask(ctx context.Context, taskID string) (*models.TaskQueue, error)
	CreateTask(ctx context.Context, t models.TaskQueue) (*models.TaskQueue, error)
}

type DBTaskQueueRepo struct {
	db *sql.DB
}

func NewDBTaskQueueRepo(db *sql.DB) *DBTaskQueueRepo {
	return &DBTaskQueueRepo{db: db}
}

func (r *DBTaskQueueRepo) GetTask(ctx context.Context, taskID string) (*models.TaskQueue, error) {
	var task models.TaskQueue
	query := "SELECT id, heading_id, status, response, cost, created_at, formatted_prompt, article_id, prompt_id FROM tasks_queue WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID,
		&task.HeadingID,
		&task.Status,
		&task.Response,
		&task.Cost,
		&task.CreatedAt,
		&task.FormattedPrompt,
		&task.ArticleID,
		&task.PromptID,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &task, err
}

func (r *DBTaskQueueRepo) CreateTask(ctx context.Context, t models.TaskQueue) (*models.TaskQueue, error) {
	query := "INSERT INTO tasks_queue (id, heading_id, status, response, cost, formatted_prompt, article_id, prompt_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.HeadingID,
		t.Status,
		t.Response,
		t.Cost,
		t.FormattedPrompt,
		t.ArticleID,
		t.PromptID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &t, nil
}
