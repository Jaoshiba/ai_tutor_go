package entities

import (
	"time"
)

type ModuleDataModel struct {
	ModuleId   string    `json:"moduleId" db:"moduleid"`
	ModuleName string    `json:"moduleName" db:"modulename"`
	RoadmapId  string    `json:"roadmapId" db:"roadmapid"`
	UserId     string    `json:"userId" db:"userid"`
	CreatedAt  time.Time `json:"createdAt" db:"createdat"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updatedat"`
}
