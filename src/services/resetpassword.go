package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"

	"log"

	"github.com/gofiber/fiber/v2"
)

type ResetPassword struct {
	ResetPasswordRepo repositories.IResetPassword
	UserRepo          repositories.IUsersRepository
}

type IResetPasswordService interface {
	CreateResetRequest(ctx *fiber.Ctx, email string) error
	ResetPassword(ctx *fiber.Ctx, token, newpassword string) error
}

func NewResetPasswordService(resetpasswordRepository repositories.IResetPassword, userRepository repositories.IUsersRepository) IResetPasswordService {
	if resetpasswordRepository == nil {
		log.Fatal("nil repo")
	}
	return &ResetPassword{
		ResetPasswordRepo: resetpasswordRepository,
		UserRepo:          userRepository,
	}
}

func GenerateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *ResetPassword) CreateResetRequest(ctx *fiber.Ctx, email string) error {
	fmt.Println("Email is : ", email)
	user, err := r.UserRepo.GetUserInfo(email)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user of this email doesn't exist")
	}

	// สร้าง token
	token, err := GenerateToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// สร้าง entity สำหรับเก็บ DB
	resetReq := &entities.ResetPassword{
		Id:        uuid.NewString(),
		UserId:    user.UserID.String(),
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsReset:   false,
		IsAvailable: true,
	}

	// เก็บลง DB
	err = r.ResetPasswordRepo.InsertNewResetPassword(context.Background(), resetReq)
	if err != nil {
		return fmt.Errorf("failed to save reset request: %w", err)
	}

	resetLink := fmt.Sprintf(os.Getenv("FRONTEND_URL")+"/resetpassword?token=%s", token)

	// สร้าง email body ในรูปแบบ HTML แบบเรียบง่าย
	emailBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body>
		<div style="font-family: Arial, sans-serif;">
			<h2>Password Reset Request</h2>
			<p>Hello %s,</p>
			<p>We received a request to reset your password. Click the button below to reset it:</p>
			<p>
				<a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
			</p>
			<p>This link will expire in 24 hours.</p>
			<p>If you did not request this, please ignore this email.</p>
			<br>
			<p style="font-size: 12px; color: #888;">This is an automated email. Please do not reply.</p>
		</div>
	</body>
	</html>
	`, user.Email, resetLink)

	// ส่ง email (placeholder)
	err = SendEmail(
		"smtp.gmail.com",
		587,
		email,
		"Reset Your Password",
		emailBody,
	)
	if err != nil {
		fmt.Println("failed to send reset email")
	}

	return nil
}

func (r *ResetPassword) ResetPassword(ctx *fiber.Ctx, token, newpassword string) error {
	// 1) ดึงข้อมูล token
	rec, err := r.ResetPasswordRepo.GetResetPasswordByToken(context.Background(), token)
	if err != nil {
		return fmt.Errorf("failed to get reset token: %w", err)
	}
	if rec == nil {
		return fmt.Errorf("invalid reset token")
	}

	// 2) ตรวจสถานะ token
	now := time.Now()
	if rec.IsReset {
		return fmt.Errorf("this reset token has already been used")
	}
	if !rec.IsAvailable {
		return fmt.Errorf("this reset token is not available")
	}
	if now.After(rec.ExpiresAt) {
		return fmt.Errorf("this reset token has expired")
	}

	// 3) แฮชรหัสผ่านใหม่
	hashed, err := HashPassword(newpassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// 4) แปลง user_id เป็น uuid.UUID
	// userUUID, err := uuid.Parse(rec.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id format in reset record: %w", err)
	}

	// 5) อัปเดตรหัสผ่านผู้ใช้
	if err := r.UserRepo.UpdateUserPassword(rec.UserId, hashed); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	// 6) ปิด token นี้ว่าใช้งานแล้ว (is_reset=true, is_available=false)
	if err := r.ResetPasswordRepo.UpdateResetPasswordStatus(context.Background(), rec.Id, true); err != nil {
		return fmt.Errorf("failed to mark reset token as used: %w", err)
	}

	// 7) ปิด token อื่นๆ ของผู้ใช้นี้ทั้งหมดเพื่อความปลอดภัย (กันนำอันเก่ามาใช้ซ้ำ)
	if err := r.ResetPasswordRepo.DeactivateAllByUserID(context.Background(), rec.UserId); err != nil {
		// ไม่ถึงกับ fail flow reset ทั้งหมด แต่ log ไว้
		fmt.Println("warning: failed to deactivate other reset tokens:", err)
	}

	return nil
}
