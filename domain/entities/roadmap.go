package entities

import (
	"time"
)

type Roadmap struct {
	RoadmapID   string    `json:"roadmap_id" db:"roadmapid"`
	RoadmapName string    `json:"roadmap_name" db:"roadmapname"`
	UserId      string    `json:"user_id" db:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
