package db 

import (
	"context"
	"database/sql"
	"fmt"
	"caravagio-api-golang/internal/app/models"
)

type ArticleRepo interface {
    GetArticle(ctx context.Context, articleID string) (models.Article, error)
}

type DBArticleRepo struct {
    db *sql.DB
}

func (s *DBArticleRepo) GetArticle(ctx context.Context, articleID string) (models.Article, error) {
	// get fist row from articles table

	var article models.Article

	err := s.db.QueryRowContext(ctx, "SELECT article_id, user_id, language, main_keywords, urls, status, keywords, heading_data, parsed_prompt, created_at, total_words, cost, html_content FROM articles WHERE article_id = ?", articleID).Scan(
		&article.ArticleID,
		&article.UserID,
		&article.Language,
		&article.MainKeywords,
		&article.URLs,
		&article.Status,
		&article.Keywords,
		&article.HeadingData,
		&article.ParsedPrompt,
		&article.CreatedAt,
		&article.TotalWords,
		&article.Cost,
		&article.HTMLContent,
	)

	if err != nil {
		fmt.Println(err)
		return article, err
	}

	return article, nil

}

func NewDBArticleRepo(db *sql.DB) *DBArticleRepo {
    return &DBArticleRepo{db: db}
}