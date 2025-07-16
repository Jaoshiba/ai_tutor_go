package entities

import "time"

type ExamDataModel struct {
	ExamId    string              `json:"examId" bson:"examId,omitempty"`
	ChapterId string              `json:"chapterId" bson:"chapterId,omitempty"`
	Questions  []QuestionDataModel `json:"question" bson:"question,omitempty"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type QuestionDataModel struct {
	Question string   `json:"question" bson:"question,omitempty"`
	Options  []string `json:"options" bson:"options,omitempty"`
	Answer   string   `json:"answer" bson:"answer,omitempty"`
}
