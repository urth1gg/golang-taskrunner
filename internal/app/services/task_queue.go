package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"log"
)

type TaskQueueService struct {
	models.TaskQueueService
}

func NewTaskQueueService(repo models.TaskQueueRepo) *TaskQueueService {
	return &TaskQueueService{models.TaskQueueService{Db: repo}}
}

// GetTask retrieves a specific task based on its ID.
func (s *TaskQueueService) GetTask(ctx context.Context, taskID string) (models.TaskQueue, error) {
	task, err := s.Db.GetTask(ctx, taskID)
	if err != nil {
		log.Printf("Failed to get task with ID %s: %v", taskID, err)
		return models.TaskQueue{}, err
	}
	return *task, nil
}

// CreateTask adds a new task to the queue.
func (s *TaskQueueService) CreateTask(ctx context.Context, task models.TaskQueue) (models.TaskQueue, error) {
	newTask, err := s.Db.CreateTask(ctx, task)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		return models.TaskQueue{}, err
	}
	return *newTask, nil
}