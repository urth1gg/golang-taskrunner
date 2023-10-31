package models

type Node struct {
	ID          string   `json:"id"`
	Tag         int      `json:"tag"`
	Text        string   `json:"text"`
	Level       int      `json:"level"`
	Length      int      `json:"length"`
	Children    []Node   `json:"children"`
	Expanded    bool     `json:"expanded"`
	Keywords    string   `json:"keywords"`
	Response    string   `json:"response"`
	MoreInfo    string   `json:"more_info"`
	PromptID    string  `json:"prompt_id"` // Assuming this can be null, hence using a pointer
	IsCompleted bool     `json:"is_completed"`
}