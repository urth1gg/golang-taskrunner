package db 

import (
	"context"
	"database/sql"
	"fmt"
	"caravagio-api-golang/internal/app/models"
	"encoding/json"
	"errors"
)

type ArticleRepo interface {
    GetArticle(ctx context.Context, articleID string) (models.Article, error)
	UpdateArticle(ctx context.Context, article *models.Article) (models.Article, error)
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

func (s *DBArticleRepo) UpdateArticle(ctx context.Context, article *models.Article) (int, error) {

	headingData, err := json.MarshalIndent(article.HeadingData, "", " ")
	if err != nil {
		fmt.Println("JSON ERR")
		fmt.Println(err)
		return 0, err
	}

	result, err := s.db.ExecContext(ctx, "UPDATE articles SET heading_data = ? WHERE article_id = ?", string(headingData), article.ArticleID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.New("no rows updated")
	}

	// Populate updatedArticle if needed, perhaps by querying the updated row
	return int(rowsAffected), nil
}

func NewDBArticleRepo(db *sql.DB) *DBArticleRepo {
    return &DBArticleRepo{db: db}
}