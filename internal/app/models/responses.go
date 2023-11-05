package models

type ArticleBody struct {
	Data struct {
		Expanded         bool   `json:"expanded"`
		ID               string `json:"id"`
		IsCompleted      bool   `json:"is_completed"`
		Keywords         string `json:"keywords"`
		Length           int    `json:"length"`
		Level            int    `json:"level"`
		MoreInfo         string `json:"more_info"`
		PromptID         string `json:"prompt_id"`
		Response         string `json:"response"`
		SettingsExpanded bool   `json:"settingsExpanded"`
		Tag              int    `json:"tag"`
		Text             string `json:"text"`
		Title            string `json:"title"`
		Children         []Node `json:"children"`
	} `json:"data"`
	Regenerate bool `json:"regenerate"`
	Continue   bool `json:"continue"`
}
