package entities

import (
	"time"
)

type QuestionDataModel struct {
	QuestionId string    `json:"questionId" db:"questionId,omitempty"`
	ExamId     string    `json:"examId" db:"examid,omitempty"`
	Type       string    `json:"type" db:"type,omitempty"`
	Question   string    `json:"question" db:"question,omitempty"`
	Optoions   string    `json:"options" db:"options,omitempty"`
	Answer     string    `json:"answer" db:"answer,omitempty"`
	CreateAt   time.Time `json:"createAt" db:"createat,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updatedat,omitempty"`
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

type QuestionFromGemini struct {
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Options  []string `json:"options,omitempty"`
	Answer   []string `json:"answer"`
}

type QuestionRequest struct {
	ExamId      string `json:"examId" db:"examId,omitempty"`
	Content     string `json:"content" db:"content,omitempty"`
	Difficulty  string `json:"difficulty" db:"difficulty,omitempty"`
	QuestionNum int    `json:"questionNum" db:"questionNum,omitempty"`
}
