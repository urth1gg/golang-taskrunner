package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"caravagio-api-golang/internal/app/handlers/articles"
	servArticles "caravagio-api-golang/internal/app/services/articles"
	servAuth "caravagio-api-golang/internal/app/services/auth"
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/middleware/auth"
	"log"
	"fmt"
)

func main() {
	r := gin.Default()

	user := "test"
	password := "mysql"
	host := "localhost"
	port := 3306
	dbName := "caravagio"


	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}  // Allow only localhost:3000 to access the API
	config.AllowHeaders = append(config.AllowHeaders, "Body")  // Allow the Body header
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")  // Allow the Authorization header
	config.AllowCredentials = true

	r.Use(cors.New(config))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	conn, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authRepo := db.NewDBAuthRepo(conn.DB)
	authService := servAuth.NewAuthService(authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(authService)

	r.Use(authMiddleware.Middleware())

	articleRepo := db.NewDBArticleRepo(conn.DB)
	articleService := servArticles.NewArticleService(articleRepo)
	articleHandler := articleshandler.NewHandler(articleService)

	

	r.GET("/articles", articleHandler.HelloWorld)
	r.GET("/articles/:articleID", articleHandler.GetArticle)
	r.POST("/articles/:articleID/generate", articleHandler.UpdateArticle)

	// Start the server
	r.Run(":8080")
}