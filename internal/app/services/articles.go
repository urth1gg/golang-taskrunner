package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
)

type ArticleService struct {
	db           db.ArticleRepo
	taskQueueSvc TaskQueueService
}

func NewArticleService(db db.ArticleRepo, taskQueueSvc *TaskQueueService) *ArticleService {
    return &ArticleService{db: db, taskQueueSvc: *taskQueueSvc}
}

func (s *ArticleService) GetArticle(ctx context.Context, articleID string) (models.Article, error) {
	// get fist row from articles table

	article, err := s.db.GetArticle(ctx, articleID)

	if err != nil {
		fmt.Println(err)
	}

	return article, nil
}
