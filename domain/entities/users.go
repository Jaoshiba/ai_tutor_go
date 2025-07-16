package entities

import (
	"time"
)

type UserDataModel struct {
	UserID     string    `json:"user_id" bson:"user_id,omitempty"`
	Username   string    `json:"username" bson:"username,omitempty"`
	Password   string    `json:"password" bson:"password,omitempty"` // อย่าลืม hash ก่อนเก็บ
	Email      string    `json:"email" bson:"email,omitempty"`
	Phone      string    `json:"phone" bson:"phone,omitempty"`
	FullName   string    `json:"full_name" bson:"full_name,omitempty"`
	Gender     string    `json:"gender" bson:"gender,omitempty"` // เช่น "male", "female", "other"
	BirthDate  time.Time `json:"birth_date" bson:"birth_date,omitempty"`
	Address    string    `json:"address" bson:"address,omitempty"` // เช่น "Bangkok, Thailand"
	Ip         string    `json:"ip" bson:"ip,omitempty"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at,omitempty"`
}


