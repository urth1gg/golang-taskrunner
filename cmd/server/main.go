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
	"os"
)

func main() {
	r := gin.Default()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	conn, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	clientChannels := make(map[string]chan services.GptResponse)

	authRepo := db.NewDBAuthRepo(conn.DB)
	authService := services.NewAuthService(authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(authService)

	r.Use(authMiddleware.Middleware())

	variablesRepo := db.NewDBVariablesRepo(conn.DB)
	promptRepo := db.NewDBPromptRepo(conn.DB)
	taskQueueRepo := db.NewDBTaskQueueRepo(conn.DB)
	articleRepo := db.NewDBArticleRepo(conn.DB)
	settingsRepo := db.NewDBSettingsRepo(conn.DB)
	defaultPromptsRepo := db.NewDBDefaultPromptsRepo(conn.DB)

	variablesService := services.NewVariablesService(variablesRepo)
	promptService := services.NewPromptService(promptRepo, variablesService)
	taskQueueService := services.NewTaskQueueService(taskQueueRepo, promptService)
	articleService := services.NewArticleService(*articleRepo)
	settingsService := services.NewSettingsService(settingsRepo)
	openAiService := services.NewOpenAIService("", clientChannels, articleService)
	defaultPromptsService := services.NewDefaultPromptsService(defaultPromptsRepo)

	streamGptHandler := handlers.NewStreamGptHandler(authService, taskQueueService, clientChannels)
	articleHandler := handlers.NewArticleHandler(articleService, taskQueueService)
	settingsHandler := handlers.NewSettingsHandler(settingsService, defaultPromptsService)

	taskExecutor := services.NewTaskExecutor(openAiService, taskQueueService, settingsService, articleService, variablesService)
	taskExecutor.RunScheduledTaskLoader(900 * time.Millisecond)
	taskExecutor.StartWorkers(10)

	r.GET("/articles/:articleID", articleHandler.GetArticle)
	r.PATCH("/articles/:articleID", articleHandler.UpdateArticle)
	r.GET("/streamgpt/:userID", streamGptHandler.SendData)
	r.DELETE("/tasks", articleHandler.DeleteTasks)
	r.PUT("/settings/:userID/default-prompts", settingsHandler.UpdateDefaultPrompts)
	r.GET("/settings/:userID/default-prompts", settingsHandler.GetSettings)
	r.Run(":8080")
}
