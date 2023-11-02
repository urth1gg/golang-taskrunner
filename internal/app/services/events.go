package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
)

type EventsService struct {
	taskQueueSvc TaskQueueService
}

func NewEventsService(taskQueueSvc *TaskQueueService) *EventsService {
	return &EventsService{taskQueueSvc: *taskQueueSvc}
}

func (s *EventsService) GetAllCompletedTasks() (*[]models.TaskQueue, error) {
	ctx := context.Background()
	tasks, err := s.taskQueueSvc.GetAllCompletedTasks(ctx)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &tasks, nil

}
