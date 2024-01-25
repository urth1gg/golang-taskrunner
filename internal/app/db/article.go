package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// TODO: Make separate function that takes in a slice of strings and updates those fields
func (s *DBArticleRepo) UpdateArticleGeneric(ctx context.Context, article *models.Article, fieldsToUpdate []string) (int, error) {
	articleID := article.ArticleID
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Rollback the transaction in case of a panic or error
	defer tx.Rollback()

	// Build the SQL statement dynamically based on the fieldsToUpdate slice
	setClauses := []string{}
	args := []interface{}{}
	for _, field := range fieldsToUpdate {
		var value interface{}
		switch field {
		case "UserID":
			value = article.UserID
		case "Language":
			value = article.Language
		case "main_keywords":
			value = article.MainKeywords
		case "URLs":
			value = article.URLs
		case "Status":
			value = article.Status
		case "keywords":
			value = article.Keywords
		case "heading_data":
			headingData, err := json.MarshalIndent(article.HeadingData, "", " ")
			if err != nil {
				fmt.Println("JSON Marshal Error:", err)
				return 0, err
			}
			value = string(headingData)
		case "ParsedPrompt":
			value = article.ParsedPrompt
		case "TotalWords":
			value = article.TotalWords
		case "Cost":
			value = article.Cost
		case "HTMLContent":
			value = article.HTMLContent
		case "IsCompleted":
			value = article.IsCompleted
		case "meta_description":
			value = article.MetaDescription
		default:
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, articleID)

	// Join all set clauses with commas
	setClause := strings.Join(setClauses, ", ")
	fmt.Printf("setClause: %s\n", setClause)
	query := fmt.Sprintf("UPDATE articles SET %s WHERE article_id = ?", setClause)

	// Execute the update query
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
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

	return int(rowsAffected), nil
}

func (s *DBArticleRepo) CreateArticle(ctx context.Context, article *models.Article) (int, error) {

	headingData, err := json.MarshalIndent(article.HeadingData.Data, "", " ")
	if err != nil {
		fmt.Println("JSON ERR")
		fmt.Println(err)
		return 0, err
	}

	result, err := s.db.ExecContext(ctx, "INSERT INTO articles (article_id, user_id, language, main_keywords, urls, status, keywords, heading_data, parsed_prompt, created_at, total_words, cost, html_content, meta_description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", article.ArticleID, article.UserID, article.Language, article.MainKeywords, article.URLs, article.Status, article.Keywords, string(headingData), article.ParsedPrompt, article.CreatedAt, article.TotalWords, article.Cost, article.HTMLContent, article.MetaDescription)

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

	return int(rowsAffected), nil
}

func NewDBArticleRepo(db *sql.DB) *DBArticleRepo {
	return &DBArticleRepo{db: db}
}

// 	+------------------+-----------------------------------+------+-----+---------+-------+
// | Field            | Type                              | Null | Key | Default | Extra |
// +------------------+-----------------------------------+------+-----+---------+-------+
// | article_id       | varchar(255)                      | NO   | PRI | NULL    |       |
// | user_id          | varchar(255)                      | YES  | MUL | NULL    |       |
// | language         | varchar(255)                      | YES  |     | NULL    |       |
// | main_keywords    | varchar(255)                      | YES  |     | NULL    |       |
// | urls             | text                              | YES  |     | NULL    |       |
// | status           | enum('active','inactive','draft') | YES  |     | draft   |       |
// | keywords         | text                              | YES  |     | NULL    |       |
// | heading_data     | json                              | YES  |     | NULL    |       |
// | parsed_prompt    | text                              | YES  |     | NULL    |       |
// | created_at       | datetime                          | YES  |     | NULL    |       |
// | total_words      | int                               | YES  |     | 0       |       |
// | cost             | decimal(18,10)                    | YES  |     | NULL    |       |
// | html_content     | text                              | YES  |     | NULL    |       |
// | meta_description | text                              | YES  |     | NULL    |       |
// +------------------+-----------------------------------+------+-----+---------+-------+
