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
	MetaTaskStatusPending      = "meta_pending"
	TaskStatusPending          = "pending"
	TaskStatusCompleted        = "completed"
	TaskStatusFailed           = "failed"
	TaskStatusRetrying         = "retrying"
	TaskStatusInProgress       = "in_progress"
	TaskStatusCompletedAndSent = "completed_and_sent"
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

// delete task

func (s *TaskQueueService) DeleteTask(ctx context.Context, task models.TaskQueue) error {
	err := s.db.DeleteTask(ctx, task)
	if err != nil {
		log.Printf("Failed to delete task: %v", err)
		return err
	}
	return nil
}

func (s *TaskQueueService) MarkTasksAsCompletedAndSent(ctx context.Context, tasks []models.TaskQueue) error {
	for _, task := range tasks {
		task.Status = TaskStatusCompletedAndSent
		_, err := s.UpdateTask(ctx, task)
		if err != nil {
			log.Printf("Failed to update task: %v", err)
			return err
		}
	}
	return nil
}

func (s *TaskQueueService) CreateTasksFromArticle(ctx context.Context, article models.Article) ([]models.TaskQueue, error) {

	tasks := []models.TaskQueue{}

	if !article.HeadingData.Data[0].IsCompleted {
		log.Println("Creating task for H1")
		t := models.TaskQueue{
			ID:              uuid.New().String(),
			ArticleID:       article.ArticleID,
			Status:          TaskStatusPending,
			HeadingID:       article.HeadingData.Data[0].ID,
			Response:        sql.NullString{String: "", Valid: false},
			Cost:            sql.NullFloat64{Float64: 0, Valid: false},
			FormattedPrompt: sql.NullString{String: "", Valid: false},
			PromptID:        article.HeadingData.Data[0].PromptID,
			GptModel:        "",
			MaxTokens:       article.HeadingData.Data[0].Length,
		}

		prompt, err := s.PromptService.GetPrompt(ctx, t.PromptID)

		t.GptModel = prompt.GPTModel.String

		if err != nil {
			log.Printf("Failed to get prompt: %v", err)
			return nil, err
		}

		formattedPrompt := s.PromptService.GenerateFormattedPromptWithAllVariablesH1(&prompt, &article)

		log.Println("Formatted Prompt")

		if err != nil {
			log.Printf("Failed to generate formatted prompt: %v", err)
			return nil, err
		}

		t.FormattedPrompt.String = formattedPrompt
		t.FormattedPrompt.Valid = true
		s.CreateTask(ctx, t)
	}

	for _, header := range article.HeadingData.Data[0].Children {

		if header.Level == 2 {
			log.Println("Creating task for H2")

			if !header.IsCompleted {
				log.Println("H2 is completed, skipping")
				prompt, err := s.PromptService.GetPrompt(ctx, header.PromptID)

				if err != nil {
					log.Printf("Failed to get prompt: %v", err)
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
					GptModel:        "",
					MaxTokens:       header.Length,
				}
				t.GptModel = prompt.GPTModel.String

				formattedPrompt, err := s.PromptService.GenerateFormattedPromptWithAllVariables(&prompt, &header, &article)

				if err != nil {
					log.Printf("Failed to generate formatted prompt: %v", err)
					return nil, err
				}

				t.FormattedPrompt.String = formattedPrompt
				t.FormattedPrompt.Valid = true

				tasks = append(tasks, t)
			}

			if len(header.Children) > 0 {
				for _, subHeader := range header.Children {
					log.Println("Creating task for H3")

					if subHeader.IsCompleted {
						log.Println("H3 is completed, skipping")
						continue
					}

					prompt, err := s.PromptService.GetPrompt(ctx, subHeader.PromptID)

					if err != nil {
						log.Printf("Failed to get prompt: %v", err)
						return nil, err
					}

					t := models.TaskQueue{
						ID:              uuid.New().String(),
						ArticleID:       article.ArticleID,
						Status:          TaskStatusPending,
						HeadingID:       subHeader.ID,
						Response:        sql.NullString{String: "", Valid: false},
						Cost:            sql.NullFloat64{Float64: 0, Valid: false},
						FormattedPrompt: sql.NullString{String: "", Valid: false},
						PromptID:        subHeader.PromptID,
						GptModel:        "",
						MaxTokens:       subHeader.Length,
					}

					t.GptModel = prompt.GPTModel.String

					formattedPrompt, err := s.PromptService.GenerateFormattedPromptWithAllVariables(&prompt, &subHeader, &article)

					if err != nil {
						log.Printf("Failed to generate formatted prompt: %v", err)
						return nil, err
					}

					t.FormattedPrompt.String = formattedPrompt
					t.FormattedPrompt.Valid = true

					tasks = append(tasks, t)
				}
			}

		}
	}

	for _, task := range tasks {
		s.CreateTask(ctx, task)
	}

	return nil, nil
}

func (s *TaskQueueService) CreateContinueTasksFromArticle(ctx context.Context, article models.Article) ([]models.TaskQueue, error) {

	tasks := []models.TaskQueue{}

	if !article.HeadingData.Data[0].IsCompleted {
		log.Println("Creating task for H1")
		t := models.TaskQueue{
			ID:                 uuid.New().String(),
			ArticleID:          article.ArticleID,
			Status:             TaskStatusPending,
			HeadingID:          article.HeadingData.Data[0].ID,
			Response:           sql.NullString{String: "", Valid: false},
			Cost:               sql.NullFloat64{Float64: 0, Valid: false},
			FormattedPrompt:    sql.NullString{String: "", Valid: false},
			PromptID:           article.HeadingData.Data[0].PromptID,
			GptModel:           "",
			ContinueGenerating: true,
			MaxTokens:          article.HeadingData.Data[0].Length,
		}

		prompt, err := s.PromptService.GetPrompt(ctx, t.PromptID)

		if err != nil {
			log.Printf("Failed to get prompt: %v", err)
			return nil, err
		}

		t.GptModel = prompt.GPTModel.String

		formattedPrompt := s.PromptService.GenerateFormattedPromptWithAllVariablesH1(&prompt, &article)

		log.Println("Formatted Prompt")

		if err != nil {
			log.Printf("Failed to generate formatted prompt: %v", err)
			return nil, err
		}

		t.FormattedPrompt.String = formattedPrompt
		t.FormattedPrompt.Valid = true
		s.CreateTask(ctx, t)
	}

	for _, header := range article.HeadingData.Data[0].Children {

		if header.Level == 2 {
			log.Println("Creating task for H2")

			if !header.IsCompleted {
				log.Println("H2 is completed, skipping")
				prompt, err := s.PromptService.GetPrompt(ctx, header.PromptID)

				if err != nil {
					log.Printf("Failed to get prompt: %v", err)
					return nil, err
				}

				t := models.TaskQueue{
					ID:                 uuid.New().String(),
					ArticleID:          article.ArticleID,
					Status:             TaskStatusPending,
					HeadingID:          header.ID,
					Response:           sql.NullString{String: "", Valid: false},
					Cost:               sql.NullFloat64{Float64: 0, Valid: false},
					FormattedPrompt:    sql.NullString{String: "", Valid: false},
					PromptID:           header.PromptID,
					GptModel:           "",
					ContinueGenerating: true,
					MaxTokens:          header.Length,
				}

				t.GptModel = prompt.GPTModel.String

				formattedPrompt, err := s.PromptService.GenerateFormattedPromptWithAllVariables(&prompt, &header, &article)

				if err != nil {
					log.Printf("Failed to generate formatted prompt: %v", err)
					return nil, err
				}

				t.FormattedPrompt.String = formattedPrompt
				t.FormattedPrompt.Valid = true

				tasks = append(tasks, t)
			}

			if len(header.Children) > 0 {
				for _, subHeader := range header.Children {
					log.Println("Creating task for H3")

					if subHeader.IsCompleted {
						log.Println("H3 is completed, skipping")
						continue
					}

					prompt, err := s.PromptService.GetPrompt(ctx, subHeader.PromptID)

					if err != nil {
						log.Printf("Failed to get prompt: %v", err)
						return nil, err
					}

					t := models.TaskQueue{
						ID:                 uuid.New().String(),
						ArticleID:          article.ArticleID,
						Status:             TaskStatusPending,
						HeadingID:          subHeader.ID,
						Response:           sql.NullString{String: "", Valid: false},
						Cost:               sql.NullFloat64{Float64: 0, Valid: false},
						FormattedPrompt:    sql.NullString{String: "", Valid: false},
						PromptID:           subHeader.PromptID,
						GptModel:           "",
						ContinueGenerating: true,
						MaxTokens:          subHeader.Length,
					}

					t.GptModel = prompt.GPTModel.String

					formattedPrompt, err := s.PromptService.GenerateFormattedPromptWithAllVariables(&prompt, &subHeader, &article)

					if err != nil {
						log.Printf("Failed to generate formatted prompt: %v", err)
						return nil, err
					}

					t.FormattedPrompt.String = formattedPrompt
					t.FormattedPrompt.Valid = true

					tasks = append(tasks, t)
				}
			}

		}
	}

	for _, task := range tasks {
		s.CreateTask(ctx, task)
	}

	return nil, nil
}

func (s *TaskQueueService) GetAllCompletedTasks(ctx context.Context) ([]models.TaskQueue, error) {
	tasks, err := s.db.GetAllCompletedTasks(ctx)
	if err != nil {
		log.Printf("Failed to get all tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (s *TaskQueueService) AddTasksToHistory(ctx context.Context, tasks []models.TaskQueue) error {
	err := s.db.AddTasksToHistory(ctx, tasks)

	if err != nil {
		log.Printf("Failed to add tasks to history: %v", err)
		return err
	}

	for _, task := range tasks {
		err := s.DeleteTask(ctx, task)

		if err != nil {
			log.Printf("Failed to delete task: %v", err)
			return err
		}
	}

	return nil
}

func (s *TaskQueueService) GetTaskFromHistoryByHeadingId(ctx context.Context, headingID string) (*models.TaskQueue, error) {
	task, err := s.db.GetTaskFromHistoryByHeadingId(ctx, headingID)

	if err != nil {
		log.Printf("Failed to get task from history: %v", err)
		return nil, err
	}

	return task, nil
}

func (s *TaskQueueService) CreateMetaDescriptionTask(ctx context.Context, article *models.Article, metaDescription *models.Node) (models.TaskQueue, error) {
	log.Println("Creating task for meta description")

	prompt, err := s.PromptService.GetPrompt(ctx, metaDescription.PromptID)

	if err != nil {
		log.Printf("Failed to get prompt: %v", err)
		return models.TaskQueue{}, err
	}

	t := models.TaskQueue{
		ID:              uuid.New().String(),
		ArticleID:       article.ArticleID,
		Status:          MetaTaskStatusPending,
		HeadingID:       metaDescription.ID,
		Response:        sql.NullString{String: "", Valid: false},
		Cost:            sql.NullFloat64{Float64: 0, Valid: false},
		FormattedPrompt: sql.NullString{String: "", Valid: false},
		PromptID:        metaDescription.PromptID,
		GptModel:        "",
		MaxTokens:       metaDescription.Length,
	}

	t.GptModel = prompt.GPTModel.String
	generatedPrompt, err := s.PromptService.GenerateFormattedPromptWithAllVariables(&prompt, metaDescription, article)

	if err != nil {
		log.Printf("Failed to generate formatted prompt: %v", err)
		return models.TaskQueue{}, err
	}

	t.FormattedPrompt.String = generatedPrompt
	t.FormattedPrompt.Valid = true

	s.CreateTask(ctx, t)

	return t, nil
}

func (s *TaskQueueService) DeleteTasks(ctx context.Context) {
	err := s.db.DeleteTasks(ctx)

	if err != nil {
		log.Printf("Failed to delete tasks: %v", err)
	}
}
