package entities

import (
	"time"
)

type FileDataModel struct {
	FileId   string    `json:"file_id" bson:"file_id,omitempty"`
	FileName string    `json:"filename" bson:"filename,omitempty"`
	UploadAt time.Time `json:"uploadAt" bson:"UploadAt,omitempty"`
	OwnerId  string    `json:"owner_id" bson:"owner_id,omitempty"`
}
