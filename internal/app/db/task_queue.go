package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type TaskQueueRepo interface {
	GetTask(ctx context.Context, taskID string) (*models.TaskQueue, error)
	CreateTask(ctx context.Context, t models.TaskQueue) (*models.TaskQueue, error)
	GetAllPendingTasks(ctx context.Context) ([]models.TaskQueue, error)
}

type DBTaskQueueRepo struct {
	db *sql.DB
}

func NewDBTaskQueueRepo(db *sql.DB) *DBTaskQueueRepo {
	return &DBTaskQueueRepo{db: db}
}

func (r *DBTaskQueueRepo) GetTask(ctx context.Context, taskID string) (*models.TaskQueue, error) {
	var task models.TaskQueue
	query := "SELECT id, heading_id, status, response, cost, created_at, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens FROM tasks_queue WHERE id = ?"
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
		&task.GptModel,
		&task.ContinueGenerating,
		&task.MaxTokens,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &task, err
}

func (r *DBTaskQueueRepo) CreateTask(ctx context.Context, t models.TaskQueue) (*models.TaskQueue, error) {
	query := "INSERT INTO tasks_queue (id, heading_id, status, response, cost, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.HeadingID,
		t.Status,
		t.Response,
		t.Cost,
		t.FormattedPrompt,
		t.ArticleID,
		t.PromptID,
		t.GptModel,
		t.ContinueGenerating,
		t.MaxTokens,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &t, nil
}

func (r *DBTaskQueueRepo) GetAllPendingTasks(ctx context.Context) ([]models.TaskQueue, error) {
	var tasks []models.TaskQueue
	query := "SELECT id, heading_id, status, response, cost, created_at, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens FROM tasks_queue WHERE status = 'pending' OR status = 'meta_pending'"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task models.TaskQueue
		err := rows.Scan(
			&task.ID,
			&task.HeadingID,
			&task.Status,
			&task.Response,
			&task.Cost,
			&task.CreatedAt,
			&task.FormattedPrompt,
			&task.ArticleID,
			&task.PromptID,
			&task.GptModel,
			&task.ContinueGenerating,
			&task.MaxTokens,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *DBTaskQueueRepo) GetAllCompletedTasks(ctx context.Context) ([]models.TaskQueue, error) {
	var tasks []models.TaskQueue
	query := "SELECT id, heading_id, status, response, cost, created_at, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens FROM tasks_queue WHERE status = 'completed'"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task models.TaskQueue
		err := rows.Scan(
			&task.ID,
			&task.HeadingID,
			&task.Status,
			&task.Response,
			&task.Cost,
			&task.CreatedAt,
			&task.FormattedPrompt,
			&task.ArticleID,
			&task.PromptID,
			&task.GptModel,
			&task.ContinueGenerating,
			&task.MaxTokens,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *DBTaskQueueRepo) UpdateTask(ctx context.Context, task models.TaskQueue) (*models.TaskQueue, error) {
	query := "UPDATE tasks_queue SET status = ?, response = ?, cost = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, task.Status, task.Response, task.Cost, task.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	query = "SELECT status, response, cost FROM tasks_queue WHERE id = ?"

	task2 := models.TaskQueue{}
	err = r.db.QueryRowContext(ctx, query, task.ID).Scan(
		&task2.Status,
		&task2.Response,
		&task2.Cost,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task2, nil
}

// If the logic gets more complicated, we can move this to a separate repo
func (r *DBTaskQueueRepo) AddTasksToHistory(ctx context.Context, tasks []models.TaskQueue) error {
	query := "INSERT INTO tasks_queue_history (id, heading_id, status, response, cost, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	for _, task := range tasks {
		_, err := r.db.ExecContext(ctx, query,
			uuid.New().String(),
			task.HeadingID,
			task.Status,
			task.Response,
			task.Cost,
			task.FormattedPrompt,
			task.ArticleID,
			task.PromptID,
			task.GptModel,
			task.ContinueGenerating,
			task.MaxTokens,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (r *DBTaskQueueRepo) GetTaskFromHistoryByHeadingId(ctx context.Context, headingID string) (*models.TaskQueue, error) {
	var task models.TaskQueue
	query := "SELECT id, heading_id, status, response, cost, created_at, formatted_prompt, article_id, prompt_id, gpt_model, continue_generating, max_tokens FROM tasks_queue_history WHERE heading_id = ? ORDER BY created_at DESC LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, headingID).Scan(
		&task.ID,
		&task.HeadingID,
		&task.Status,
		&task.Response,
		&task.Cost,
		&task.CreatedAt,
		&task.FormattedPrompt,
		&task.ArticleID,
		&task.PromptID,
		&task.GptModel,
		&task.ContinueGenerating,
		&task.MaxTokens,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &task, err
}

func (r *DBTaskQueueRepo) DeleteTask(ctx context.Context, task models.TaskQueue) error {
	query := "DELETE FROM tasks_queue WHERE heading_id = ?"
	_, err := r.db.ExecContext(ctx, query, task.HeadingID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
