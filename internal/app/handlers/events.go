package handlers

import (
	"caravagio-api-golang/internal/app/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EventsHandler struct {
	EventsService    *services.EventsService
	AuthService      *services.AuthService
	TaskQueueService *services.TaskQueueService
}

func (h *EventsHandler) SendData(c *gin.Context) {
	authValue, exists := c.Get("Authorization")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authHeader, ok := authValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	apiKey, err := h.AuthService.ValidateAPIKey(c, authHeader)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("userID")

	if apiKey.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	w := c.Writer
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	for {
		ctx := context.Background()
		tasks, err := h.EventsService.GetAllCompletedTasks()

		if len(*tasks) == 0 {
			time.Sleep(555 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Println(err)
		}

		h.TaskQueueService.MarkTasksAsCompletedAndSent(ctx, *tasks)
		err = h.TaskQueueService.AddTasksToHistory(ctx, *tasks)

		if err != nil {
			fmt.Println(err)
		}

		if len(*tasks) > 0 {

			data, err := json.Marshal(gin.H{"tasks": tasks})
			if err != nil {
				fmt.Println("Failed to marshal tasks:", err)
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", data)

		}

		flusher.Flush()
		time.Sleep(555 * time.Millisecond)
	}

}

func NewEventsHandler(eventsService *services.EventsService, authService *services.AuthService, taskQueueService *services.TaskQueueService) *EventsHandler {
	return &EventsHandler{EventsService: eventsService, AuthService: authService, TaskQueueService: taskQueueService}
}
