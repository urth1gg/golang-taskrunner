package handlers

import (
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StreamGptHandler struct {
	EventsService    *services.EventsService
	AuthService      *services.AuthService
	TaskQueueService *services.TaskQueueService
	Response         *chan services.GptResponse
}

func (h *StreamGptHandler) SendData(c *gin.Context) {
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
	//w.Header().Set("Access-Control-Allow-Origin", "http://143.110.157.129:3000")
	w.Header().Set("Transfer-Encoding", "chunked")

	clientCtx := c.Request.Context()

	fmt.Println("Client connected to stream")
	for {
		select {
		case <-clientCtx.Done():
			fmt.Println("Client disconnected")
			return
		case response := <-*h.Response:
			json, err := json.Marshal(response)

			if err != nil {
				fmt.Println("Failed to marshal tasks:", err)
				fmt.Fprintf(w, "event: %s\n", "message")
				fmt.Fprintf(w, "data: %s\n\n", "{\"error\": \"Failed to marshal tasks\"}")
				continue
			}

			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: %s\n\n", json)

			flusher.Flush()
		}
	}

}

func NewStreamGptHandler(eventsService *services.EventsService, authService *services.AuthService, taskQueueService *services.TaskQueueService, responseChannel *chan services.GptResponse) *StreamGptHandler {
	return &StreamGptHandler{EventsService: eventsService, AuthService: authService, TaskQueueService: taskQueueService, Response: responseChannel}
}
