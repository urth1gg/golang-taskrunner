package handlers

import (
	"caravagio-api-golang/internal/app/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var connectedUsers = make(map[string]interface{})
var mutex = sync.Mutex{}

type StreamGptHandler struct {
	AuthService      *services.AuthService
	TaskQueueService *services.TaskQueueService
	ClientChannels   map[string]chan services.GptResponse
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

	if connectedUsers[userID] == true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User already connected"})
		return
	}

	w := c.Writer
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Set("Access-Control-Allow-Origin", "http://143.110.157.129:3000")
	w.Header().Set("Transfer-Encoding", "chunked")

	clientCtx := c.Request.Context()

	connectedUsers[userID] = true

	fmt.Println("Client connected to stream")

	clientChannel := make(chan services.GptResponse)
	h.ClientChannels[userID] = clientChannel

	heartbeatInterval := 30 * time.Second
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-clientCtx.Done():
			fmt.Println("Client disconnected")
			flusher.Flush()
			ticker.Stop()

			mutex.Lock()
			if ch, ok := h.ClientChannels[userID]; ok {
				close(ch)
				delete(h.ClientChannels, userID)
				delete(connectedUsers, userID)
			}
			mutex.Unlock()

			return
		case response := <-h.ClientChannels[userID]:
			json, err := json.Marshal(response)

			if err != nil {
				fmt.Println("Failed to marshal tasks:", err)
				fmt.Fprintf(w, "event: %s\n", "message")
				fmt.Fprintf(w, "data: %s\n\n", "{\"error\": \"Failed to marshal tasks\"}")
				flusher.Flush()
				continue
			}

			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: %s\n\n", json)

			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, "event: %s\n", "heartbeat")
			fmt.Fprintf(w, "data: %s\n\n", "{}")
			flusher.Flush()
		}
	}

}

func NewStreamGptHandler(authService *services.AuthService, taskQueueService *services.TaskQueueService, clientChannels map[string]chan services.GptResponse) *StreamGptHandler {
	return &StreamGptHandler{AuthService: authService, TaskQueueService: taskQueueService, ClientChannels: clientChannels}
}
