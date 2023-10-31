package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
)

type ArticleService struct {
	db           db.DBArticleRepo
	taskQueueSvc TaskQueueService
}

func NewArticleService(db db.DBArticleRepo, taskQueueSvc *TaskQueueService) *ArticleService {
    return &ArticleService{db: db, taskQueueSvc: *taskQueueSvc}
}

func (s *ArticleService) GetArticle(ctx context.Context, articleID string) (models.Article, error) {

	article, err := s.db.GetArticle(ctx, articleID)

	if err != nil {
		fmt.Println(err)
	}

	return article, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, article *models.Article) (int, error) {

	affectedRows, err := s.db.UpdateArticle(ctx, article)

	if err != nil {
		fmt.Println(err)
	}

	return affectedRows, nil
}
