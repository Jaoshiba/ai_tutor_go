package entities

import "time"

type ExamDataModel struct {
	ExamId      string    `json:"examId" db:"examId,omitempty"`
	ChapterId   string    `json:"chapterId" db:"chapterId,omitempty"`
	PassScore   int       `json:"passScore" db:"passScore,omitempty"`
	QuestionNum int       `json:"questionNum" db:"questionNum,omitempty"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedAt,omitempty"`
}
