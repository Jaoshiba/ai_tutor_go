package entities

import (
	"time"
)

type CourseDataModel struct {
	CourseID    string    `json:"course_id" db:"courseid"`
	CourseName  string    `json:"course_name" db:"title"`
	Description string    `json:"description" db:"description"`
	Confirmed   bool      `json:"confirmed" db:"confirmed"`
	UserId      string    `json:"user_id" db:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseRequestBody struct {
	CourseName  string `json:"courseName" db:"title"`
	Description string `json:"description" db:"description"`
	Confirmed   bool   `json:"confirmed" db:"confirmed"`
}

type CourseGeminiResponse struct {
	Modules []GenModule `json:"modules"`
}

type GenModule struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
