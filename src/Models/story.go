package Models

import "time"

type Stories struct {
	StoryId int64 `json:"id"`
	Title   string `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AllStoryAPIResponse struct {
	Limit int         `json:"limit"`
	Offset int        `json:"offset"`
	Count int         `json:"count"`
	Results []Stories `json:"results"`
}

type Story struct {
	Stories
	Paragraphs []Paragraph `json:"paragraphs"`
}