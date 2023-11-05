package models

import (
	"database/sql"
	"github.com/google/uuid"
)

type TaskQueue struct {
	ID                 string          `db:"id"`
	HeadingID          string          `db:"heading_id"`
	Status             string          `db:"status"`
	Response           sql.NullString  `db:"response"`
	Cost               sql.NullFloat64 `db:"cost"`
	CreatedAt          NullTime        `db:"created_at"`
	FormattedPrompt    sql.NullString  `db:"formatted_prompt"`
	ArticleID          string          `db:"article_id"`
	PromptID           string          `db:"prompt_id"`
	GptModel           string          `db:"gpt_model"`
	ContinueGenerating bool            `db:"continue_generating"`
}

func NewTaskQueue(headingID, status, response, formattedPrompt, articleID, promptID string, cost float64, gptModel string, continueGenerating bool) TaskQueue {
	return TaskQueue{
		ID:                 uuid.New().String(),
		HeadingID:          headingID,
		Status:             status,
		Response:           sql.NullString{String: response, Valid: true},
		Cost:               sql.NullFloat64{Float64: cost, Valid: true},
		FormattedPrompt:    sql.NullString{String: formattedPrompt, Valid: true},
		ArticleID:          articleID,
		PromptID:           promptID,
		GptModel:           gptModel,
		ContinueGenerating: continueGenerating,
	}
}
