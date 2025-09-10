package entities

import "time"

type ExamDataModel struct {
	ExamId      string    `json:"examId" db:"examid,omitempty"`
	ChapterId   string    `json:"chapterId" db:"chapterid,omitempty"`
	PassScore   int       `json:"passScore" db:"passscore,omitempty"`
	QuestionNum int       `json:"questionNum" db:"questionnum,omitempty"`
	Difficulty  string    `json:"difficulty" db:"difficulty,omitempty"`
	CreatedAt   time.Time `json:"createdAt" db:"createdat,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedat,omitempty"`
}

type ExamRequest struct {
	ChapterId   string `json:"chapterId" db:"chapterId,omitempty"`
	Content     string `json:"content" db:"content,omitempty"`
	Difficulty  string `json:"difficulty" db:"difficulty,omitempty"`
	QuestionNum int    `json:"questionNum" db:"questionNum,omitempty"`
}
