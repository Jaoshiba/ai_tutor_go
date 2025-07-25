package entities

import (
	"time"
)

type RoadmapDataModel struct {
	RoadmapID   string    `json:"roadmap_id" db:"roadmapid"`
	RoadmapName string    `json:"roadmap_name" db:"roadmapname"`
	Description string    `json:"description" db:"description"`
	Confirmed   bool      `json:"confirmed" db:"confirmed"`
	UserId      string    `json:"user_id" db:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoadmapRequestBody struct {
	RoadmapName string `json:"roadmapName" db:"roadmapname"`
	Description string `json:"description" db:"description"`
	Confirmed   bool   `json:"confirmed" db:"confirmed"`
}

type RoadmapGeminiResponse struct {
	Modules []GenModule `json:"modules"`
}

type GenModule struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
