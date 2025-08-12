package entities

import (
	"time"
)

type ChapterDataModel struct {
	ChapterId      string    `json:"chapterId" db:"chapterid"`
	ChapterName    string    `json:"chapterName" db:"chaptername"`
	UserID         string    `json:"userId" db:"userid"`
	CourseId       string    `json:"courseId" db:"courseid"`
	ChapterContent string    `json:"chapterContent" db:"chaptercontent"`
	IsFinished     bool      `json:"isFinished" db:"isfinished"`
	CreateAt       time.Time `json:"createAt" db:"createat"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updatedat"`
	ModuleId       string    `json:"moduleid" db:"moduleid"`
}

type ResponseChapter struct {
	ChapterName string `json:"chapterName" db:"chaptername,omitempty"`
	Content     string `json:"content" db:"content,omitempty"`
}

type GeminiResponse struct {
	Message  string            `json:"message" db:"message,omitempty"`
	Chapters []ResponseChapter `json:"chapters" db:"chapters,omitempty"`
}
