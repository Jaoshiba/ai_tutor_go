package entities

import (
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`  // จะเก็บ user info + token
	Error   string      `json:"error,omitempty"` // ข้อความ error ถ้ามี
}

type UserInfoResponse struct {
	UserID uuid.UUID `json:"userid"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}