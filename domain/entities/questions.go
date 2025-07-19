package entities

import (
	"time"
)

type QuestionDataModel struct {
	QuestionId string    `json:"questionId" db:"questionId,omitempty"`
	Question   string    `json:"question" db:"question,omitempty"`
	ChapterId  string    `json:"chapterId" db:"chapterId,omitempty"`
	CorrAnsId  string    `json:"corrAnsId" db:"corrAnsId,omitempty"`
	CreateAt   time.Time `json:"createAt" db:"createAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updatedAt,omitempty"`
}
