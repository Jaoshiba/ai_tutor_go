package entities

import (
	"time"
)

type ChapterDataModel struct {
	ChapterId   string    `json:"chapterId" bson:"chapterId,omitempty"`
	ChapterName string    `json:"chapterName" bson:"chapterName,omitempty"`
	Content     string    `json:"content" bson:"content,omitempty"`
	CreateAt    time.Time `json:"createAt" bson:"createAt,omitempty"`
}

type ResponseChapter struct {
	ChapterName string `json:"chapterName" bson:"chapterName,omitempty"`
	Content     string `json:"content" bson:"content,omitempty"`
}

type GeminiResponse struct {
	Message  string            `json:"message" bson:"message,omitempty"`
	Chapters []ResponseChapter `json:"chapters" bson:"chapters,omitempty"`
}
