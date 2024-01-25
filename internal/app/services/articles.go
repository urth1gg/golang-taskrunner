package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
)

type ArticleService struct {
	db db.DBArticleRepo
}

func NewArticleService(db db.DBArticleRepo) *ArticleService {
	return &ArticleService{db: db}
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

func (s *ArticleService) UpdateArticleGeneric(ctx context.Context, article *models.Article, fieldsToUpdate []string) (int, error) {

	affectedRows, err := s.db.UpdateArticleGeneric(ctx, article, fieldsToUpdate)

	if err != nil {
		fmt.Println(err)
	}

	return affectedRows, nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *models.Article) (int, error) {

	affectedRows, err := s.db.CreateArticle(ctx, article)

	if err != nil {
		fmt.Println(err)
	}

	return affectedRows, nil
}
