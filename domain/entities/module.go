package entities

import (
	"time"
)

type ModuleDataModel struct {
	ModuleId   string             `json:"moduleId" bson:"moduleId,omitempty"`
	ModuleName string             `json:"moduleName" bson:"moduleName,omitempty"`
	RoadmapId  string             `json:"roadmapId" bson:"roadmapId,omitempty"`
	Chapters   []ChapterDataModel `json:"chapters" bson:"chapters,omitempty"`
	Exam       []ExamDataModel    `json:"exam" bson:"exam,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}
