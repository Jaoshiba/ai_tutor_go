package entities

import (
	"time"
)

type RefDataModel struct {
	RefId    string    `json:"refId" db:"refid"`
	ModuleId string    `json:"moduleId" db:"moduleid"`
	Title    string    `json:"title" db:"title"`
	Link     string    `json:"link" db:"link"`
	Content  string    `json:"content" db:"content"`
	SearchAt time.Time `json:"searchAt" db:"searchat"`
}
