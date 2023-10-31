package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"log"
	"time"
)

type TaskExecutor struct {
	TaskQueue        chan Task
	RetryQueue       chan Task
	MaxRetries       int
	TaskQueueService *TaskQueueService
	OpenAIService    *OpenAIService
	RetryDelay       time.Duration
}

type Task struct {
	Data      models.TaskQueue
	Retries   int
	LastError error
}

func NewTaskExecutor(openAiService *OpenAIService, taskQueueService *TaskQueueService) *TaskExecutor {
	return &TaskExecutor{
		TaskQueue:        make(chan Task, 100),
		RetryQueue:       make(chan Task, 100),
		MaxRetries:       3,
		RetryDelay:       1 * time.Minute,
		TaskQueueService: taskQueueService,
	}
}

func (te *TaskExecutor) LoadPendingTasks(ctx context.Context) {
	tasks, err := te.TaskQueueService.GetAllPendingTasks(ctx)
	if err != nil {
		log.Printf("Error fetching pending tasks: %v", err)
		return
	}
	for _, task := range tasks {
		te.TaskQueue <- Task{Data: task}
	}
}

func (te *TaskExecutor) StartWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go te.worker()
	}
	go te.retryWorker()
}

func (te *TaskExecutor) worker() {
	for task := range te.TaskQueue {
		err := te.processTask(task.Data)
		if err != nil {
			task.Retries++
			task.LastError = err
			if task.Retries <= te.MaxRetries {
				te.RetryQueue <- task
			} else {
				log.Printf("Task failed after %d attempts: %v", te.MaxRetries, err)
				task.Data.Status = TaskStatusFailed
				_, updateErr := te.TaskQueueService.UpdateTask(context.Background(), task.Data)
				if updateErr != nil {
					log.Printf("Error updating task status: %v", updateErr)
				}
			}
		} else {
			task.Data.Status = TaskStatusCompleted
			_, updateErr := te.TaskQueueService.UpdateTask(context.Background(), task.Data)
			if updateErr != nil {
				log.Printf("Error updating task status: %v", updateErr)
			}
		}
	}
}

func (te *TaskExecutor) retryWorker() {
	for task := range te.RetryQueue {
		time.Sleep(te.RetryDelay)
		te.TaskQueue <- task
	}
}

func (te *TaskExecutor) processTask(taskData models.TaskQueue) error {
	ctx := context.Background()

	log.Println("PROCESSING TASK")
	if taskData.GptModel == "gpt-3.5" {
		resp, err := te.OpenAIService.UseGPT3_5(ctx, taskData.FormattedPrompt.String)

		if err != nil {
			log.Println(err)
			return err
		}

		taskData.Response.String = resp
		taskData.Response.Valid = true

		_, err = te.TaskQueueService.UpdateTask(ctx, taskData)

		if err != nil {
			log.Println(err)
			return err
		}

		log.Println(resp)
	}

	return nil
}

func (te *TaskExecutor) AddTask(data models.TaskQueue) {
	te.TaskQueue <- Task{Data: data}
}

func (te *TaskExecutor) RunScheduledTaskLoader(interval time.Duration) {
	ticker := time.NewTicker(interval)

	log.Println("RUNNING")
	go func() {
		for range ticker.C {
			ctx := context.Background()
			te.LoadPendingTasks(ctx)
		}
	}()
}
