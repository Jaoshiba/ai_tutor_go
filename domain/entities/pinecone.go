package entities

import (
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

type UpsertBodyChapter struct {
	Namespace string             `json:"namespace"`
	Vectors   []*pinecone.Vector `json:"vectors"`
}

type MetadataChapter struct {
	ChapterId   string `json:"chapterid"`
	Chaptername string `json:"chaptername"`
	UserId      string `json:"userid"`
	CourseId    string `json:"courseid"`
}

func (m *MetadataChapter) ToMap() map[string]interface{} {
	if m == nil { // ป้องกัน nil pointer dereference
		return nil
	}
	return map[string]interface{}{
		"chapterid":   m.ChapterId,
		"chaptername": m.Chaptername,
		"userid":      m.UserId,
		"courseid":    m.CourseId,
	}
}
