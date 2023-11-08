package models

type Article struct {
	ArticleID       string      `sql:"article_id" json:"article_id"`
	UserID          string      `sql:"user_id" json:"user_id"`
	Language        string      `sql:"language" json:"language"`
	MainKeywords    string      `sql:"main_keywords" json:"main_keywords"`
	URLs            *string     `sql:"urls" json:"urls"` // If URLs are comma or newline separated
	Status          *string     `sql:"status" json:"status"`
	Keywords        string      `sql:"keywords" json:"keywords"`         // If Keywords are comma or newline separated
	HeadingData     HeadingData `sql:"heading_data" json:"heading_data"` // Assuming JSON is marshaled as string
	ParsedPrompt    *string     `sql:"parsed_prompt" json:"parsed_prompt"`
	CreatedAt       NullTime    `sql:"created_at" json:"created_at"`
	TotalWords      *int        `sql:"total_words" json:"total_words"`
	Cost            float64     `sql:"cost" json:"cost"` // Assuming cost is decimal with up to 10 decimal places
	HTMLContent     *string     `sql:"html_content" json:"html_content"`
	MetaDescription string      `sql:"meta_description" json:"meta_description"`
	IsCompleted     bool        `json:"is_completed"`
	MoreInfo        string      `json:"more_info"`
	Length          int         `json:"length"`
}
