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

type ExamRequest struct {
	ModuleId    string `json:"module" db:"module,omitempty"`
	Content     string `json:"content" db:"content,omitempty"`
	Difficulty  string `json:"difficulty" db:"difficulty,omitempty"`
	QuestionNum int    `json:"questionNum" db:"questionNum,omitempty"`
}
