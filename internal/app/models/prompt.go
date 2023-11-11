package models

import (
	"database/sql"
)

type Prompt struct {
	PromptID              string          `db:"prompt_id"`
	UserID                sql.NullString  `db:"user_id"`
	Name                  sql.NullString  `db:"name"`
	Description           sql.NullString  `db:"description"`
	TextArea              sql.NullString  `db:"text_area"`
	GPTModel              sql.NullString  `db:"gpt_model"`
	Temperature           sql.NullFloat64 `db:"temperature"`
	MaxLength             sql.NullInt64   `db:"max_length"`
	TopP                  sql.NullFloat64 `db:"top_p"`
	FrequencyPenalty      sql.NullFloat64 `db:"frequency_penalty"`
	PresencePenalty       sql.NullFloat64 `db:"presence_penalty"`
	CreatedAt             NullTime        `db:"created_at"`
	LevelRequiredToAccess string          `db:"level_required_to_access"`
}
