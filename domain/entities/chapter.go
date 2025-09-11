package entities

import (
	"time"
)

type ChapterDataModel struct {
	ChapterId      string    `json:"chapterId" db:"chapterid"`
	ChapterName    string    `json:"chapterName" db:"chaptername"`
	ChapterContent string    `json:"chapterContent" db:"chaptercontent"`
	CreateAt       time.Time `json:"createAt" db:"createat"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updatedat"`
	ModuleId       string    `json:"moduleid" db:"moduleid"`
	Description    string    `json:"description" db:"description"`
	Index          int       `json:"index" db:"index"`
	Ispassed       bool      `json:"ispassed" db:"ispassed"`
}

type ResponseChapter struct {
	ChapterName string `json:"chapterName" db:"chaptername,omitempty"`
	Content     string `json:"content" db:"content,omitempty"`
}

type GeminiResponse struct {
	Message  string            `json:"message" db:"message,omitempty"`
	Chapters []ResponseChapter `json:"chapters" db:"chapters,omitempty"`
}
