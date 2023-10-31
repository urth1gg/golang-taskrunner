package main

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/handlers"
	"caravagio-api-golang/internal/app/middleware"
	"caravagio-api-golang/internal/app/services"
	// "caravagio-api-golang/internal/app/models"
	// "context"
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
	config.AllowOrigins = []string{"http://localhost:3000"}            // Allow only localhost:3000 to access the API
	config.AllowHeaders = append(config.AllowHeaders, "Body")          // Allow the Body header
	config.AllowHeaders = append(config.AllowHeaders, "Authorization") // Allow the Authorization header
	config.AllowCredentials = true

	r.Use(cors.New(config))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	conn, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	promptRepo := db.NewDBPromptRepo(conn.DB)
	promptService := services.NewPromptService(promptRepo)

	taskQueueRepo := db.NewDBTaskQueueRepo(conn.DB)
	taskQueueService := services.NewTaskQueueService(taskQueueRepo, promptService)

	authRepo := db.NewDBAuthRepo(conn.DB)
	authService := services.NewAuthService(authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(authService)

	r.Use(authMiddleware.Middleware())

	articleRepo := db.NewDBArticleRepo(conn.DB)
	articleService := services.NewArticleService(*articleRepo, taskQueueService)
	articleHandler := articleshandler.NewHandler(articleService)

	openAiService := services.NewOpenAIService("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")
	taskExecutor := services.NewTaskExecutor(openAiService, taskQueueService)
	taskExecutor.RunScheduledTaskLoader(5 * time.Minute) // Run every 5 minutes
	taskExecutor.StartWorkers(10)                        // Start 10 workers

	// ctx := context.Background()
	// taskQueue := models.NewTaskQueue("123", "done", "response", "formattedPrompt", "ff753636-0098-44ae-bd6e-aceb72ee5efc", "831a62f3-7504-11ee-b232-00155db8d1fa", 0.0)

	// _, err = taskQueueService.CreateTask(ctx, taskQueue)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// ctx := context.Background()
	// article, err := articleService.GetArticle(ctx, "81754eaa-caf6-4f16-92ec-e57aff969863")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// taskQueueService.CreateTasksFromArticle(ctx, article)
	// promptText := `
	// 	Write engaging and interesting introduction to article about the topic {h1_title}. Content should be unconventional and attract the attention of the reader. When the reader reads the introduction, he should want to read the whole article. Refer to the reader, but don't write words like "hi", "welcome", "hello" etc. Below you can find all headings that should be in the article, as well as keywords - so you know what the whole article is about. You can (but you don't have to) use the keywords and headings while writing the introduction.

	// 	All headings:
	// 	{all_header}

	// 	All keywords:
	// 	{keywords}
	// `

	// prompt := models.Prompt{
	// 	PromptID: "831a62f3-7504-11ee-b232-00155db8d1fa",
	// 	TextArea: sql.NullString{
	// 		String: promptText,
	// 		Valid:  true,
	// 	},
	// }

	// formattedText, err := articleService.GenerateFormattedPromptH1Intro(&prompt, &article)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(formattedText)

	r.GET("/articles", articleHandler.HelloWorld)
	r.GET("/articles/:articleID", articleHandler.GetArticle)
	r.PATCH("/articles/:articleID", articleHandler.UpdateArticle)

	//r.POST("/articles/:articleID/generate", articleHandler.UpdateArticle)

	// Start the server
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
