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

type ChoiceOption struct {
	QuestionId   string   `json:"questionId" db:"questionId,omitempty"`
	QuestionType string   `json:"type" db:"type,omitempty"`
	QuestionText string   `json:"question" db:"question,omitempty"`
	Options      []string `json:"options" db:"options,omitempty"`
	Answer       string   `json:"answer" db:"answer,omitempty"`
}

type FillInBlank struct {
	QuestionId   string `json:"questionId" db:"questionId,omitempty"`
	QuestionType string `json:"type" db:"type,omitempty"`
	QuestionText string `json:"question" db:"question,omitempty"`
	Answer       string `json:"answer" db:"answer,omitempty"`
}

type Ordering struct {
	QuestionId   string   `json:"questionId" db:"questionId,omitempty"`
	QuestionType string   `json:"type" db:"type,omitempty"`
	QuestionText string   `json:"question" db:"question,omitempty"`
	Options      []string `json:"options" db:"options,omitempty"`
	Answer       []string `json:"answer" db:"answer,omitempty"`
}
