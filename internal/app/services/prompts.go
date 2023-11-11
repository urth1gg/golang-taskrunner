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

func (s *PromptService) GenerateFormattedPromptWithAllVariablesH1(prompt *models.Prompt, article *models.Article) string {
	if !prompt.TextArea.Valid {
		// Handle the case where the TextArea is null.
		// You can return an empty string, an error, or some default value.
		return ""
	}

	mainHeading := article.MainKeywords
	headersData := s.GenerateAllHeadersText(article)
	keywords := article.Keywords
	moreInfo := article.MoreInfo

	fmt.Printf("%s\n", keywords)
	fmt.Printf("%s\n", moreInfo)
	// TODO: Remove the duplicate string replace calls such as {Keywords} and {keywords} and {more_info} and {additional_info}.
	// We do not know the exact information that was in the previous version.

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", mainHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.Replace(formattedText, "{keywords}", keywords, -1)
	formattedText = strings.Replace(formattedText, "{Keywords}", keywords, -1)
	formattedText = strings.Replace(formattedText, "{more_info}", moreInfo, -1)
	formattedText = strings.Replace(formattedText, "{additional_info}", moreInfo, -1)

	formattedText = strings.TrimSpace(formattedText)

	return formattedText
}

func (s *PromptService) GenerateAllHeadersText(article *models.Article) string {
	headers := []string{}

	// Assuming MainKeywords is a slice of strings and you want to include it as headers
	headers = append(headers, article.MainKeywords)

	// Recursively append all headers
	s.appendHeaders(&headers, article.HeadingData.Data[0].Children)

	// Join all headers with two newlines
	headersText := strings.Join(headers, "\n\n")

	return headersText
}

// Recursive function to append headers
func (s *PromptService) appendHeaders(headers *[]string, nodes []models.Node) {
	for _, node := range nodes {
		*headers = append(*headers, node.Text)
		if len(node.Children) > 0 {
			s.appendHeaders(headers, node.Children)
		}
	}
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

// TODO: Refactor this
func (s *PromptService) findHeading(headingData []models.Node, node *models.Node, findPrev bool) (string, bool) {
	var lastHeadingAtLevel string
	found := false

	for _, heading := range headingData {
		// If we've found our node, return the last heading at the same level
		if found && heading.Level == node.Level && findPrev {
			return lastHeadingAtLevel, true
		}

		// If this is the node, mark it as found
		if heading.Text == node.Text && heading.Level == node.Level {
			found = true
			// If we're looking for the next heading, continue the search
			if !findPrev {
				continue
			}
		}

		// If we're looking for the previous heading and we've found the node,
		// the last heading at the level is our answer
		if found && findPrev {
			return lastHeadingAtLevel, true
		}

		// Keep track of the last heading at this level
		if heading.Level == node.Level {
			lastHeadingAtLevel = heading.Text
		}

		// Recursively search in children
		if len(heading.Children) > 0 {
			if result, foundInChild := s.findHeading(heading.Children, node, findPrev); foundInChild {
				return result, true
			}
		}
	}

	// If we're looking for the next heading and we've found the node,
	// but there's no next heading at the same level, return an empty string
	if found && !findPrev {
		return "", true
	}

	return "", false
}

func (s *PromptService) findNextHeader(headingData []models.Node, targetNode *models.Node, foundTarget *bool) (string, bool) {
	for _, heading := range headingData {
		if *foundTarget {
			if heading.Level == targetNode.Level {
				return heading.Text, true
			}
		} else if heading.Text == targetNode.Text && heading.Level == targetNode.Level {
			*foundTarget = true
		}

		if len(heading.Children) > 0 {
			if nextHeader, found := s.findNextHeader(heading.Children, targetNode, foundTarget); found {
				return nextHeader, true
			}
		}
	}

	return "", false
}

func (s *PromptService) GeneratePrevHeader(node *models.Node, article *models.Article) string {
	prevHeading, _ := s.findHeading(article.HeadingData.Data[0].Children, node, true)
	return prevHeading
}

func (s *PromptService) GenerateNextHeader(node *models.Node, article *models.Article) string {
	foundTarget := false
	nextHeader, _ := s.findNextHeader(article.HeadingData.Data[0].Children, node, &foundTarget)
	return nextHeader
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
	parentHeader := s.GenerateParentHeader(node, article)
	keywords := node.Keywords
	moreInfo := node.MoreInfo

	fmt.Printf("%s\n", keywords)
	fmt.Printf("%s\n", moreInfo)
	// TODO: Remove the duplicate string replace calls such as {Keywords} and {keywords} and {more_info} and {additional_info}.
	// We do not know the exact information that was in the previous version.

	formattedText := strings.Replace(prompt.TextArea.String, "{h1_title}", mainHeading, -1)
	formattedText = strings.Replace(formattedText, "{h2_title}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{all_header}", headersData, -1)
	formattedText = strings.Replace(formattedText, "{current_header}", nodeHeading, -1)
	formattedText = strings.Replace(formattedText, "{previous_header}", prevHeader, -1)
	formattedText = strings.Replace(formattedText, "{next_header}", nextHeader, -1)
	formattedText = strings.Replace(formattedText, "{keywords}", keywords, -1)
	formattedText = strings.Replace(formattedText, "{Keywords}", keywords, -1)
	formattedText = strings.Replace(formattedText, "{parent_header}", parentHeader, -1)
	formattedText = strings.Replace(formattedText, "{more_info}", moreInfo, -1)
	formattedText = strings.Replace(formattedText, "{additional_info}", moreInfo, -1)

	formattedText = strings.TrimSpace(formattedText)

	return formattedText, nil
}

func (s *PromptService) GetAllAvailablePrompts(ctx context.Context, levelRequiredToAccess string) ([]models.Prompt, error) {
	prompts, err := s.db.GetAllAvailablePrompts(ctx, levelRequiredToAccess)

	if err != nil {
		fmt.Println(err)
		return prompts, err
	}

	return prompts, nil
}
