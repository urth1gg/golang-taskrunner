package services

import (
	"errors"
	"strings"
	// "caravagio-api-golang/internal/app/db"
	"context"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"log"
)

type GptResponse struct {
	HeadingID string
	Response  string
}

type OpenAIService struct {
	client         *openai.Client
	ClientChannels map[string]chan GptResponse
	ArticleService *ArticleService
}

// Initialize a new OpenAIService with your OpenAI API key.
func NewOpenAIService(apiKey string, clientChannels map[string]chan GptResponse, articleService *ArticleService) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{client: client, ClientChannels: clientChannels, ArticleService: articleService}
}

func (s *OpenAIService) SetOpenAIKey(apiKey string) {
	s.client = openai.NewClient(apiKey)
}

func (s *OpenAIService) UseGPT3_5(ctx context.Context, inputText string, headingID string, maxTokens int, model string, taskID string, articleID string) (string, error) {
	log.Println("Using GPT3.5")

	if strings.TrimSpace(inputText) == "" {
		return "", errors.New("input text cannot be empty")
	}

	message := openai.ChatCompletionMessage{
		Content: inputText,
		Role:    "system",
	}
	request := openai.ChatCompletionRequest{
		Messages:  []openai.ChatCompletionMessage{message},
		Model:     model,
		MaxTokens: maxTokens,
	}

	response, err := s.client.CreateChatCompletionStream(ctx, request)

	if err != nil {
		return "", err
	}

	article, err := s.ArticleService.GetArticle(ctx, articleID)

	if err != nil {
		return "", err
	}

	var result string

	// Use a for loop to receive messages from the stream
	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, return an error
			return "", ctx.Err()
		default:
			// Try to receive a message from the stream
			msg, err := response.Recv()

			if shouldResponseStreamsBeCancelled[taskID] == true {
				log.Println("Cancelling stream")
				delete(shouldResponseStreamsBeCancelled, taskID)
				goto END
			}

			if err == io.EOF {
				// If no more messages are coming through the stream, break the loop
				goto END
			}
			if err != nil {
				// For any other error, return it
				return "", err
			}

			// Process the message (assuming the response has a field 'Text' for the result)
			result += msg.Choices[0].Delta.Content
			chanData := GptResponse{
				HeadingID: headingID,
				Response:  msg.Choices[0].Delta.Content,
			}
			// *s.Response <- chanData

			if s.ClientChannels[article.UserID] != nil {
				log.Println("Sending response to client")
				s.ClientChannels[article.UserID] <- chanData
			}
		}
	}

END:
	log.Println("Result: ", result)
	return result, nil
}

func (s *OpenAIService) UseGPT4(ctx context.Context, inputText string, headingID string, maxTokens int, model string, taskID string, articleID string) (string, error) {

	if strings.TrimSpace(inputText) == "" {
		return "", errors.New("input text cannot be empty")
	}

	message := openai.ChatCompletionMessage{
		Content: inputText,
		Role:    "system",
	}
	request := openai.ChatCompletionRequest{
		Messages:  []openai.ChatCompletionMessage{message},
		Model:     model,
		MaxTokens: maxTokens,
	}

	response, err := s.client.CreateChatCompletionStream(ctx, request)

	if err != nil {
		return "", err
	}

	article, err := s.ArticleService.GetArticle(ctx, articleID)

	if err != nil {
		return "", err
	}

	var result string

	// Use a for loop to receive messages from the stream
	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, return an error
			return "", ctx.Err()
		default:
			// Try to receive a message from the stream
			msg, err := response.Recv()

			if shouldResponseStreamsBeCancelled[taskID] == true {
				log.Println("Cancelling stream")
				delete(shouldResponseStreamsBeCancelled, taskID)
				goto END
			}
			if err == io.EOF {
				// If no more messages are coming through the stream, break the loop
				goto END
			}
			if err != nil {
				// For any other error, return it
				return "", err
			}

			// Process the message (assuming the response has a field 'Text' for the result)
			result += msg.Choices[0].Delta.Content
			chanData := GptResponse{
				HeadingID: headingID,
				Response:  msg.Choices[0].Delta.Content,
			}

			if s.ClientChannels[article.UserID] != nil {
				log.Println("Sending response to client")
				s.ClientChannels[article.UserID] <- chanData
			}
		}
	}

END:
	log.Println("Result: ", result)
	return result, nil
}
