package services

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
)

func HandleError(err error) {
	_, file, line, ok := runtime.Caller(1) // 1 level up in the stack

	if ok {
		fmt.Printf("error: %v, file: %v, line: %d\n", err, file, line)
	}
}

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
		HandleError(err)
	}

	articleUserId := article.UserID

	settings, err := te.SettingsService.GetSetting(ctx, articleUserId)

	te.OpenAIService.SetOpenAIKey(settings.APIKey.String)

	if err != nil {
		HandleError(err)
	}

	resp := ""

	if taskData.ContinueGenerating {
		prependToTheStartOfThePrompt := "Continue in the same language the text below without duplicating the content:\n\n"

		task, err := te.TaskQueueService.GetTaskFromHistoryByHeadingId(ctx, taskData.HeadingID)

		if err != nil {
			HandleError(err)
		}

		prevResponse := task.Response.String
		//prevPrompt := taskData.FormattedPrompt.String
		taskData.FormattedPrompt.String = prependToTheStartOfThePrompt + prevResponse

		fmt.Println("Prompt2")
		fmt.Printf("%s", taskData.FormattedPrompt.String)
	}

	if taskData.GptModel == "gpt-4-1106-preview" {
		resp, err = te.OpenAIService.UseGPT4(ctx, taskData.FormattedPrompt.String, taskData.HeadingID, taskData.MaxTokens, "gpt-4-1106-preview")
	} else if taskData.GptModel == "gpt-4" {
		resp, err = te.OpenAIService.UseGPT4(ctx, taskData.FormattedPrompt.String, taskData.HeadingID, taskData.MaxTokens, "gpt-4")
	} else if taskData.GptModel == "gpt-3.5-turbo" {
		resp, err = te.OpenAIService.UseGPT3_5(ctx, taskData.FormattedPrompt.String, taskData.HeadingID, taskData.MaxTokens, "gpt-3.5-turbo")
	} else if taskData.GptModel == "gpt-3.5-turbo-16k" {
		resp, err = te.OpenAIService.UseGPT3_5(ctx, taskData.FormattedPrompt.String, taskData.HeadingID, taskData.MaxTokens, "gpt-3.5-turbo")
	} else if taskData.GptModel == "gpt-3.5-turbo-1106" {
		resp, err = te.OpenAIService.UseGPT3_5(ctx, taskData.FormattedPrompt.String, taskData.HeadingID, taskData.MaxTokens, "gpt-3.5-turbo-1106")
	}

	fmt.Println(taskData.GptModel)
	fmt.Println("model")

	if err != nil {
		HandleError(err)
		//return err
	}

	if taskData.ContinueGenerating {
		resp = taskData.Response.String + resp
		taskData.Response.String = resp

	} else {
		taskData.Response.String = resp
	}

	if taskData.Status == MetaTaskStatusPending {
		article.MetaDescription = taskData.Response.String
		fields := []string{"meta_description"}

		_, err := te.ArticleService.UpdateArticleGeneric(ctx, &article, fields)

		if err != nil {
			HandleError(err)
		}
	}

	taskData.Response.Valid = true

	taskData.Status = TaskStatusCompletedAndSent

	tasks := []models.TaskQueue{taskData}

	err = te.TaskQueueService.AddTasksToHistory(ctx, tasks)

	if err != nil {
		log.Println(err)
		HandleError(err)
		//return err
	}

	// TODO: this doesn't happen due to AddTasksToHistory removing it in case of GPT-4
	// _, err = te.TaskQueueService.UpdateTask(ctx, taskData)

	if err != nil {
		HandleError(err)
		//	return err
	}

	log.Println("Going once")
	log.Println(taskData.ID)

	found := updateResponse(&article.HeadingData.Data, taskData.HeadingID, taskData.Response.String)

	if !found {
		err := fmt.Errorf("could not find heading ID %s in article %s", taskData.HeadingID, articleId)

		log.Println(err)
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
