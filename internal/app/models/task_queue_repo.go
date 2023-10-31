package models

import (
	"context"
)

type TaskQueueRepo interface {
	GetTask(ctx context.Context, taskID string) (*TaskQueue, error)
	CreateTask(ctx context.Context, t TaskQueue) (*TaskQueue, error)
	GetAllPendingTasks(ctx context.Context) ([]TaskQueue, error)
	UpdateTask(ctx context.Context, t TaskQueue) (*TaskQueue, error)
}