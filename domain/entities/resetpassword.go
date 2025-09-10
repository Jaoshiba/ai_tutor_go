package entities

import (
	"time"
)

type ResetPassword struct {
	Id string `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Token string `json:"token" db:"token"`
	IsReset bool `json:"is_reset" db:"is_reset"`
}

type ResetPasswordRequest struct {
    Email string `json:"email"`
}
