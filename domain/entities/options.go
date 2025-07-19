package entities

import (
	"time"
)

type OptionsDataModel struct {
	OptionId   string    `json:"optionId" db:"optionId,omitempty"`
	OptionText string    `json:"optionText" db:"optionText,omitempty"`
	ExamId     string    `json:"examId" db:"examId,omitempty"`
	CreateAt   time.Time `json:"createAt" db:"createAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updatedAt,omitempty"`
}
