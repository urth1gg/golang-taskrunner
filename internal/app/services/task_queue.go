package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type TaskQueueService struct {
	db            models.TaskQueueRepo
	PromptService *PromptService
}

const (
	TaskStatusPending    = "pending"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
	TaskStatusRetrying   = "retrying"
	TaskStatusInProgress = "in_progress"
)

func NewTaskQueueService(repo models.TaskQueueRepo, promptService *PromptService) *TaskQueueService {
	return &TaskQueueService{db: repo, PromptService: promptService}
}

// GetTask retrieves a specific task based on its ID.
func (s *TaskQueueService) GetTask(ctx context.Context, taskID string) (models.TaskQueue, error) {
	task, err := s.db.GetTask(ctx, taskID)
	if err != nil {
		log.Printf("Failed to get task with ID %s: %v", taskID, err)
		return models.TaskQueue{}, err
	}
	return *task, nil
}

// CreateTask adds a new task to the queue.
func (s *TaskQueueService) CreateTask(ctx context.Context, task models.TaskQueue) (models.TaskQueue, error) {
	newTask, err := s.db.CreateTask(ctx, task)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		return models.TaskQueue{}, err
	}
	return *newTask, nil
}

func (s *TaskQueueService) GetAllPendingTasks(ctx context.Context) ([]models.TaskQueue, error) {
	tasks, err := s.db.GetAllPendingTasks(ctx)
	if err != nil {
		log.Printf("Failed to get all tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (s *TaskQueueService) UpdateTask(ctx context.Context, task models.TaskQueue) (models.TaskQueue, error) {
	updatedTask, err := s.db.UpdateTask(ctx, task)
	if err != nil {
		log.Printf("Failed to update task: %v", err)
		return models.TaskQueue{}, err
	}
	return *updatedTask, nil
}

func (s *TaskQueueService) CreateTasksFromArticle(ctx context.Context, article models.Article) ([]models.TaskQueue, error) {

	tasks := []models.TaskQueue{}

	t := models.TaskQueue{
		ID:              uuid.New().String(),
		ArticleID:       article.ArticleID,
		Status:          TaskStatusPending,
		HeadingID:       article.HeadingData.Data[0].ID,
		Response:        sql.NullString{String: "", Valid: false},
		Cost:            sql.NullFloat64{Float64: 0, Valid: false},
		FormattedPrompt: sql.NullString{String: "", Valid: false},
		PromptID:        article.HeadingData.Data[0].PromptID,
		GptModel:        "gpt-3.5",
	}

	prompt, err := s.PromptService.GetPrompt(ctx, t.PromptID)

	if err != nil {
		log.Printf("Failed to get prompt: %v", err)
		return nil, err
	}

	// Generate formatted prompt

	formattedPrompt, err := s.PromptService.GenerateFormattedPromptH1Intro(&prompt, &article)

	if err != nil {
		log.Printf("Failed to generate formatted prompt: %v", err)
		return nil, err
	}

	t.FormattedPrompt.String = formattedPrompt
	t.FormattedPrompt.Valid = true

	s.CreateTask(ctx, t)

	for _, header := range article.HeadingData.Data[0].Children {

		if header.Level == 2 {
			prompt, err := s.PromptService.GetPrompt(ctx, header.PromptID)

			if err != nil {
				log.Printf("Failed to get prompt: %v", err)
				return nil, err
			}

			// Generate formatted prompt

			_, err = s.PromptService.GenerateFormattedPromptH2Intro(&prompt, &header, &article)

			if err != nil {
				log.Printf("Failed to generate formatted prompt: %v", err)
				return nil, err
			}

			t := models.TaskQueue{
				ID:              uuid.New().String(),
				ArticleID:       article.ArticleID,
				Status:          TaskStatusPending,
				HeadingID:       header.ID,
				Response:        sql.NullString{String: "", Valid: false},
				Cost:            sql.NullFloat64{Float64: 0, Valid: false},
				FormattedPrompt: sql.NullString{String: "", Valid: false},
				PromptID:        header.PromptID,
				GptModel:        "gpt-3.5",
			}

			formattedPrompt, err := s.PromptService.GenerateFormattedPromptH2Intro(&prompt, &header, &article)

			if err != nil {
				log.Printf("Failed to generate formatted prompt: %v", err)
				return nil, err
			}

			t.FormattedPrompt.String = formattedPrompt
			t.FormattedPrompt.Valid = true

			tasks = append(tasks, t)

		}
	}

	for _, task := range tasks {
		s.CreateTask(ctx, task)
	}

	return nil, nil
}
