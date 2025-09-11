package entities

import (
	"time"
)

type SearchLinks struct {
	LinkID   string    `json:"linkId" bson:"linkid"`
	ModuleID string    `json:"moduleId" bson:"moduleid"`
	Title    string    `json:"title" bson:"title"`
	Link     string    `json:"link" bson:"link"`
	Snippet  string    `json:"snippet" bson:"snippet"`
	Searchat time.Time `json:"searchat" bson:"searchat"`
}
