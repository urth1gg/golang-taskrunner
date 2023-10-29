package services

import (
	"context"
	"fmt"
	"caravagio-api-golang/internal/app/models"
	"caravagio-api-golang/internal/app/db"
)

type ArticleService struct {
	db db.ArticleService
}

func NewArticleService(db db.ArticleService) *ArticleService {
    return &ArticleService{db: db}
}

func (s *ArticleService) GetArticle(ctx context.Context, articleID string) (models.Article, error) {
	// get fist row from articles table


	article, err := s.db.GetArticle(ctx, articleID)

	if err != nil {
		fmt.Println(err)
	}

	return article, nil
}

