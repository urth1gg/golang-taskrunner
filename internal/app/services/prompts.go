package services

import (
	"caravagio-api-golang/internal/app/db"
	"caravagio-api-golang/internal/app/models"
	"context"
	"errors"
	"fmt"
	"strings"
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

func (s *PromptService) GenerateFormattedPromptH1Intro(prompt *models.Prompt, article *models.Article) string {
	if !prompt.TextArea.Valid {
		// Handle the case where the TextArea is null.
		// You can return an empty string, an error, or some default value.
		return ""
	}

	articleHeading := article.MainKeywords
	headersData := s.GenerateAllHeadersText(article)
	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", articleHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.TrimSpace(formattedText)

	return formattedText
}

func (s *PromptService) GenerateAllHeadersText(article *models.Article) string {
	HeadingData := article.HeadingData

	headers := []string{}

	headers = append(headers, article.MainKeywords)

	for _, heading := range HeadingData.Data[0].Children {
		headers = append(headers, heading.Text)
	}

	headersText := strings.Join(headers, "\n\n")

	return headersText
}

func (s *PromptService) GenerateFormattedPromptH2Intro(prompt *models.Prompt, node *models.Node, article *models.Article) (string, error) {
	if !prompt.TextArea.Valid {
		// Handle the case where the TextArea is null.
		// You can return an empty string, an error, or some default value.
		return "", errors.New("TextArea is null")
	}

	mainHeading := article.MainKeywords
	nodeHeading := node.Text
	headersData := s.GenerateAllHeadersText(article)

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", mainHeading, -1)
	formattedText = strings.Replace(formattedText, "{h2_title}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.Replace(formattedText, "{current_header}", nodeHeading, -1)
	formattedText = strings.TrimSpace(formattedText)

	return formattedText, nil
}

func (s *PromptService) GeneratePrevHeader(node *models.Node, article *models.Article) string {
	HeadingData := article.HeadingData
	prevHeading := ""

	for _, heading := range HeadingData.Data[0].Children {
		if heading.Text == node.Text {
			break
		}
		prevHeading = heading.Text
	}

	return prevHeading
}

func (s *PromptService) GenerateNextHeader(node *models.Node, article *models.Article) string {
	HeadingData := article.HeadingData
	foundCurrentNode := false

	for _, heading := range HeadingData.Data[0].Children {
		if foundCurrentNode {
			return heading.Text
		}

		if heading.Text == node.Text {
			foundCurrentNode = true
		}
	}

	return ""
}

func (s *PromptService) GenerateParentHeader(node *models.Node, article *models.Article) string {
	HeadingData := article.HeadingData

	for _, heading := range HeadingData.Data[0].Children {
		if heading.Text == node.Text {
			return HeadingData.Data[0].Text
		}

		for _, child := range heading.Children {
			if child.Text == node.Text {
				return heading.Text
			}
		}
	}

	return ""
}

func (s *PromptService) GenerateFormattedPromptWithAllVariables(prompt *models.Prompt, node *models.Node, article *models.Article) (string, error) {
	mainHeading := article.MainKeywords
	nodeHeading := node.Text
	headersData := s.GenerateAllHeadersText(article)
	prevHeader := s.GeneratePrevHeader(node, article)
	nextHeader := s.GenerateNextHeader(node, article)
	keywords := node.Keywords
	parentHeader := s.GenerateParentHeader(node, article)

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", mainHeading, -1)
	formattedText = strings.Replace(formattedText, "{h2_title}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.Replace(formattedText, "{current_header}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{previous_header}", prevHeader, -1)
	formattedText = strings.Replace(formattedText, "{next_header}", nextHeader, -1)
	formattedText = strings.Replace(formattedText, "{keywords}", keywords, -1)
	formattedText = strings.Replace(formattedText, "{parent_header}", parentHeader, -1)
	formattedText = strings.TrimSpace(formattedText)

	return formattedText, nil
}
