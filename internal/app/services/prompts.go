package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"fmt"
	"strings"
	"errors"
)

type PromptService struct {
	db db.PromptRepo
}

func NewPromptService(db db.PromptRepo) *PromptService {
	return &PromptService{db: db}
}

func (s *PromptService) GetPrompt(ctx context.Context, promptID string) (models.Prompt, error) {
	// get first row from prompts table based on promptID

	prompt, err := s.db.GetPrompt(ctx, promptID)

	if err != nil {
		fmt.Println(err)
		return models.Prompt{}, err
	}

	return *prompt, nil
}

func (s *PromptService) GenerateFormattedPromptH1Intro(prompt *models.Prompt, article *models.Article) (string, error) {
	if !prompt.TextArea.Valid {
		// Handle the case where the TextArea is null. 
		// You can return an empty string, an error, or some default value.
		return "", errors.New("TextArea is null")
	}

	articleHeading := article.MainKeywords
	headersData, err := s.GenerateAllHeadersText(article)

	if err != nil {
		fmt.Println(err)
	}

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", articleHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.TrimSpace(formattedText)

	return formattedText, nil
}

func (s *PromptService) GenerateAllHeadersText(article *models.Article) (string, error) {
	HeadingData := article.HeadingData

	headers := []string{}

	headers = append(headers, article.MainKeywords)

	for _, heading := range HeadingData.Data[0].Children {
		headers = append(headers, heading.Text)
	}

	headersText := strings.Join(headers, "\n\n")


	return headersText, nil
}

func (s *PromptService) GenerateFormattedPromptH2Intro(prompt *models.Prompt, node *models.Node, article *models.Article) (string, error){
	if !prompt.TextArea.Valid {
		// Handle the case where the TextArea is null. 
		// You can return an empty string, an error, or some default value.
		return "", errors.New("TextArea is null")
	}

	mainHeading := article.MainKeywords
	nodeHeading := node.Text 
	headersData, err := s.GenerateAllHeadersText(article)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(nodeHeading)

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", mainHeading, -1)
	formattedText = strings.Replace(formattedText, "{h2_title}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.Replace(formattedText, "{current_header}", nodeHeading, -1)
	formattedText = strings.TrimSpace(formattedText)

	return formattedText, nil
}