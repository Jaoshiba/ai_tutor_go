package entities

import (
	"time"
)

type RoadmapDataModel struct {
	RoadmapID   string    `json:"roadmap_id" db:"roadmapid"`
	RoadmapName string    `json:"roadmap_name" db:"roadmapname"`
	Description string    `json:"description" db:"description"`
	UserId      string    `json:"user_id" db:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoadmapRequestBody struct {
	
}
