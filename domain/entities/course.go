package entities

import (
	"time"
)

type CourseDataModel struct {
	CourseID    string    `json:"course_id" db:"id"`
	Title       string    `json:"Title" db:"title"`
	Description string    `json:"description" db:"description"`
	Confirmed   bool      `json:"confirmed" db:"confirmed"`
	UserId      string    `json:"user_id" db:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseRequestBody struct {
	Title       string `json:"Title" db:"title"`
	Description string `json:"Description" db:"description"`
	Confirmed   bool   `json:"confirmed" db:"confirmed"`
}

type CourseGeminiResponse struct {
	Purpose string      `json:"purpose"`
	Modules []GenModule `json:"modules"`
}

type GenModule struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type ModuleDetail struct {
	ModuleId    string          `json:"module_id"`
	ModuleName  string          `json:"module_name"`
	Description string          `json:"description"`
	Chapters    []ChapterDetail `json:"chapters"`
	// คุณอาจจะเพิ่ม CreatedAt, UpdatedAt ถ้าต้องการแสดง
}
type ChapterDetail struct {
	ChapterId      string `json:"chapter_id"`
	ChapterName    string `json:"chapter_name"`
	ChapterContent string `json:"chapter_content"`
	IsFinished     bool   `json:"is_finished"`
	// คุณอาจจะเพิ่ม CreatedAt, UpdatedAt ถ้าต้องการแสดง
}
type CourseDetailResponse struct {
	CourseID    string         `json:"course_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Confirmed   bool           `json:"confirmed"`
	Modules     []ModuleDetail `json:"modules"`
	// คุณอาจจะเพิ่ม UserId, CreatedAt, UpdatedAt ถ้าต้องการแสดง
}
