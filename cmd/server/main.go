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

/*
CREATE TABLE default_prompts (
    user_id INT NOT NULL,
    heading_name_and_position VARCHAR(255) NOT NULL,
    prompt_id INT NOT NULL,
    PRIMARY KEY (user_id, heading_name_and_position)
);
*/

/*
+--------------------------------------+---------------------------+--------------------------------------+
| user_id                              | heading_name_and_position | prompt_id                            |
+--------------------------------------+---------------------------+--------------------------------------+
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_first_h2   | c84302d4-405f-4375-bd4b-db53ef2d7842 |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_first_h3   | f5295e64-4c8f-4e09-a816-d3df8f1c4fcc |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_first_h4   | 4b3caaac-65b0-40e4-84cb-8dcec1eca10a |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_h1         | 31f25465-3224-40e9-b535-d9dc4e81b2ce |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_last_h2    | c84302d4-405f-4375-bd4b-db53ef2d7842 |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_last_h3    | 831a62f3-7504-11ee-b232-00155db8d1fa |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_last_h4    | 4b3caaac-65b0-40e4-84cb-8dcec1eca10a |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_middle_h2  | c84302d4-405f-4375-bd4b-db53ef2d7842 |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_middle_h3  | 5a1cc720-d796-43a7-8ec5-2a5bede0b28c |
| eb360df1-d5a6-467b-8bda-7ba02e114e4c | default_prompt_middle_h4  | 4b3caaac-65b0-40e4-84cb-8dcec1eca10a |

*/
