package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
	"log"
	"time"
)

type TaskExecutor struct {
	TaskQueue        chan Task
	RetryQueue       chan Task
	MaxRetries       int
	TaskQueueService *TaskQueueService
	OpenAIService    *OpenAIService
	SettingsService  *SettingsService
	ArticleService   *ArticleService
	RetryDelay       time.Duration
}

type Task struct {
	Data      models.TaskQueue
	Retries   int
	LastError error
}

func NewTaskExecutor(openAiService *OpenAIService, taskQueueService *TaskQueueService, settingsService *SettingsService, articleService *ArticleService) *TaskExecutor {
	return &TaskExecutor{
		TaskQueue:        make(chan Task, 100),
		RetryQueue:       make(chan Task, 100),
		MaxRetries:       3,
		RetryDelay:       1 * time.Minute,
		TaskQueueService: taskQueueService,
		OpenAIService:    openAiService,
		SettingsService:  settingsService,
		ArticleService:   articleService,
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
		task.Status = TaskStatusInProgress

		_, err := te.TaskQueueService.UpdateTask(ctx, task)

		if err != nil {
			fmt.Println(err)
		}

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
		}
	}
}

func (te *TaskExecutor) retryWorker() {
	for task := range te.RetryQueue {
		time.Sleep(te.RetryDelay)
		te.TaskQueue <- task
	}
}

func updateResponse(headingData *[]models.Node, headingID string, response string) bool {
	for i, node := range *headingData {
		if node.ID == headingID {
			if (*headingData)[i].IsCompleted {
				fmt.Println("Already completed")
				return false
			}
			(*headingData)[i].Response = response
			(*headingData)[i].IsCompleted = true
			return true
		}
		// If the node has children, search recursively
		if len((*headingData)[i].Children) > 0 {
			if updateResponse(&((*headingData)[i].Children), headingID, response) {
				(*headingData)[i].IsCompleted = true
				return true // Node found and updated in children
			}
		}
	}
	return false // Node not found
}

func (te *TaskExecutor) processTask(taskData models.TaskQueue) error {
	ctx := context.Background()

	articleId := taskData.ArticleID

	article, err := te.ArticleService.GetArticle(ctx, articleId)

	if err != nil {
		log.Println(err)
		return err
	}

	articleUserId := article.UserID

	settings, err := te.SettingsService.GetSetting(ctx, articleUserId)

	te.OpenAIService.SetOpenAIKey(settings.APIKey.String)

	if err != nil {
		log.Println(err)
		return err
	}

	resp := ""
	if taskData.GptModel == "gpt-4" {
		resp, err = te.OpenAIService.UseGPT4(ctx, taskData.FormattedPrompt.String, taskData.HeadingID)
	} else if taskData.GptModel == "gpt-3.5" {
		resp, err = te.OpenAIService.UseGPT3_5(ctx, taskData.FormattedPrompt.String)
	} else if taskData.GptModel == "ada-001" {
		resp, err = te.OpenAIService.UseAda(ctx, taskData.FormattedPrompt.String)
	}

	fmt.Println(taskData.GptModel)
	fmt.Println("model")
	if err != nil {
		log.Println(err)
		return err
	}

	if taskData.ContinueGenerating {
		resp = taskData.Response.String + resp
		taskData.Response.String = resp
	} else {
		taskData.Response.String = resp
	}

	taskData.Response.Valid = true
	taskData.Status = TaskStatusCompleted

	_, err = te.TaskQueueService.UpdateTask(ctx, taskData)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Going once")
	log.Println(taskData.ID)

	found := updateResponse(&article.HeadingData.Data, taskData.HeadingID, taskData.Response.String)

	if !found {
		return fmt.Errorf("could not find heading ID %s in article %s", taskData.HeadingID, articleId)
	}

	_, err = te.ArticleService.UpdateArticle(ctx, &article)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (te *TaskExecutor) AddTask(data models.TaskQueue) {
	te.TaskQueue <- Task{Data: data}
}

func (te *TaskExecutor) RunScheduledTaskLoader(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			ctx := context.Background()
			te.LoadPendingTasks(ctx)
		}
	}()
}
