package models

import (
	"context"
)

type TaskQueueRepo interface {
	GetTask(ctx context.Context, taskID string) (*TaskQueue, error)
	CreateTask(ctx context.Context, t TaskQueue) (*TaskQueue, error)
}