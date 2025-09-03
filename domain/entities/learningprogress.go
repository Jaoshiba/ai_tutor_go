
package entities

import "time"

type LearningProgressDataModel struct {
    ProcessID string    `json:"processid"`
    UserID    string    `json:"userid"`
    ModuleID  string    `json:"moduleid"`
    ChapterID string    `json:"chapterid"`
    CreatedAt time.Time `json:"createdat"`
}
