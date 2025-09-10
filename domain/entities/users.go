package entities

import (
	"time"
	"github.com/google/uuid"
)

type UserDataModel struct {
	UserID    uuid.UUID    `json:"userid" bson:"userid,omitempty"`
	Username  string    `json:"username" bson:"username,omitempty"`
	FirstName string    `json:"firstname" bson:"firstname,omitempty"`
	LastName  string    `json:"lastname" bson:"lastname,omitempty"`
	Email     string    `json:"email" bson:"email,omitempty"`
	Gender    string    `json:"gender" bson:"gender,omitempty"`
	Role      string    `json:"role" bson:"role,omitempty"` // เช่น user, admin
	DOB       string `json:"dob" bson:"dob,omitempty"`   // Date of Birth
	Password  string    `json:"password" bson:"password,omitempty"` // ต้องทำการ hash ก่อนเก็บ
	CreatedAt time.Time `json:"createdat" bson:"createdat,omitempty"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat,omitempty"`
	IsEmailVerified bool `json:"isemailverified" bson:"isemailverified,omitempty"`
}

type UserInfoModel struct {
	UserID    uuid.UUID    `json:"userid" bson:"userid,omitempty"`
	Username  string    `json:"username" bson:"username,omitempty"`
	FirstName string    `json:"firstname" bson:"firstname,omitempty"`
	LastName  string    `json:"lastname" bson:"lastname,omitempty"`
	Email     string    `json:"email" bson:"email,omitempty"`
	Gender    string    `json:"gender" bson:"gender,omitempty"`
	Role      string    `json:"role" bson:"role,omitempty"` // เช่น user, admin
	DOB       string `json:"dob" bson:"dob,omitempty"`   // Date of Birth
	CreatedAt time.Time `json:"createdat" bson:"createdat,omitempty"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat,omitempty"`
	IsEmailVerified bool `json:"isemailverified" bson:"isemailverified,omitempty"`
}

