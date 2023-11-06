package main

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/handlers"
	"caravagio-api-golang/internal/app/middleware"
	"caravagio-api-golang/internal/app/services"
	// "caravagio-api-golang/internal/app/models"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	// "database/sql"
)

func main() {
	r := gin.Default()

	user := "test"
	password := "mysql"
	host := "localhost"
	port := 3306
	dbName := "caravagio"

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://143.110.157.129:3000"} // Allow only localhost:3000 to access the API
	config.AllowHeaders = append(config.AllowHeaders, "Body")                              // Allow the Body header
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")                     // Allow the Authorization header
	config.AllowHeaders = append(config.AllowHeaders, "Access-Control-Allow-Origin")
	config.AllowHeaders = append(config.AllowHeaders, "Accept")
	config.AllowHeaders = append(config.AllowHeaders, "Cache-Control")
	config.AllowHeaders = append(config.AllowHeaders, "Transfer-Encoding")
	config.AllowCredentials = true

	r.Use(cors.New(config))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	conn, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	responseChannel := make(chan services.GptResponse)

	authRepo := db.NewDBAuthRepo(conn.DB)
	authService := services.NewAuthService(authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(authService)

	r.Use(authMiddleware.Middleware())

	promptRepo := db.NewDBPromptRepo(conn.DB)
	promptService := services.NewPromptService(promptRepo)

	taskQueueRepo := db.NewDBTaskQueueRepo(conn.DB)
	taskQueueService := services.NewTaskQueueService(taskQueueRepo, promptService)

	articleRepo := db.NewDBArticleRepo(conn.DB)
	articleService := services.NewArticleService(*articleRepo, taskQueueService)
	articleHandler := handlers.NewArticleHandler(articleService, taskQueueService)

	settingsRepo := db.NewDBSettingsRepo(conn.DB)
	settingsService := services.NewSettingsService(settingsRepo)

	openAiService := services.NewOpenAIService("", &responseChannel)
	taskExecutor := services.NewTaskExecutor(openAiService, taskQueueService, settingsService, articleService)
	taskExecutor.RunScheduledTaskLoader(900 * time.Millisecond) // Run every 5 minutes
	taskExecutor.StartWorkers(10)                               // Start 10 workers

	eventsService := services.NewEventsService(taskQueueService)
	eventsHandler := handlers.NewEventsHandler(eventsService, authService, taskQueueService)

	StreamGptHandler := handlers.NewStreamGptHandler(eventsService, authService, taskQueueService, &responseChannel)
	r.GET("/articles/:articleID", articleHandler.GetArticle)
	r.PATCH("/articles/:articleID", articleHandler.UpdateArticle)
	//r.POST("/articles/:articleID/regenerate", articleHandler.RegenerateHandler)
	r.GET("/events/:userID", eventsHandler.SendData)
	r.GET("/streamgpt/:userID", StreamGptHandler.SendData)

	r.Run(":8080")
}

// r.GET("/sse", func(c *gin.Context) {
// 	c.Writer.Header().Set("Content-Type", "text/event-stream")
// 	c.Writer.Header().Set("Cache-Control", "no-cache")
// 	c.Writer.Header().Set("Connection", "keep-alive")
// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

// 	for {
// 		// You can send any data here; in this example, we're sending the current time.
// 		data := fmt.Sprintf("data: %s\n\n", time.Now().String())

// 		// Write to the response body. This sends the data to the client.
// 		c.Writer.Write([]byte(data))

// 		// Flush the data immediately instead of buffering it.
// 		c.Writer.Flush()

// 		// Delay to simulate some data processing.
// 		time.Sleep(time.Second)
// 	}
// })
