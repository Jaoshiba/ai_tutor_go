package entities

import (
	"time"
)

type EmailVerificationModel struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Email       string    `json:"email" db:"email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	Token       string    `json:"token" db:"token"`
	IsVerify    bool      `json:"is_verify" db:"is_verify"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
