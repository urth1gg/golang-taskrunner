package services 

import (
	"errors"
	"strings"
	// "caravagio-api-golang/internal/app/db"
	openai "github.com/sashabaranov/go-openai"
	"context"
)

type OpenAIService struct {
	client *openai.Client
}

// Initialize a new OpenAIService with your OpenAI API key.
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{client: client}
}

func (s *OpenAIService) UseGPT3_5(ctx context.Context, inputText string) (string, error) {
	if strings.TrimSpace(inputText) == "" {
		return "", errors.New("input text cannot be empty")
	}

	request := openai.CompletionRequest{
		Prompt:    inputText,
		Model:     "text-davinci-003",
		MaxTokens: 150,
	}

	response, err := s.client.CreateCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}

func (s *OpenAIService) UseGPT4(ctx context.Context, inputText string) (string, error) {
	if strings.TrimSpace(inputText) == "" {
		return "", errors.New("input text cannot be empty")
	}

	request := openai.CompletionRequest{
		Prompt:    inputText,
		Model:     "gpt-4.0-turbo",
		MaxTokens: 150,
	}

	response, err := s.client.CreateCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
