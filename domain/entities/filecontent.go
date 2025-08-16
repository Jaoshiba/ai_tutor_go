package entities



type FileContentModel struct {
	Id string `json:"id" db:"id"`
	CourseId string `json:"courseid" db:"courseid"`
	Content string `json:"content" db:"content"`
	FilePath string `json:"filepath" db:"file_path"`
	CreatedAt string `json:"createat" db:"createat"`
}