package entities

import (
	"time"
)

type ModuleDataModel struct {
	ModuleId    string    `json:"moduleId" db:"moduleid"`
	ModuleName  string    `json:"moduleName" db:"modulename"`
	CourseId    string    `json:"courseId" db:"courseid"`
	UserId      string    `json:"userId" db:"userid"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"createdAt" db:"createat"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updateat"`
	Description string    `json:"description" db:"description"`
}

type ModuleDataModelGet struct {
	ModuleId    string `json:"moduleId" db:"moduleid"`
	ModuleName  string `json:"moduleName" db:"modulename"`
	CourseId    string `json:"courseId" db:"courseid"`
	UserId      string `json:"userId" db:"userid"`
	Content     string `json:"content" db:"content"`
	CreatedAt   string `json:"createdAt" db:"createat"`
	UpdatedAt   string `json:"updatedAt" db:"updateat"`
	Description string `json:"description" db:"description"`
}

// type ModuleGemini struct {
//     Title       string   `json:"title"`
//     Description string   `json:"description"`
//     Topics      []string `json:"topics"`
// }
