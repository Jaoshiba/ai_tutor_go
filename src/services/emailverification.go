package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"
)

// IEmailVerificationService defines the interface for the email verification service.
type IEmailVerificationService interface {
	CreateVerificationRequest(ctx *fiber.Ctx, email string) error
	VerifyEmail(token string) error
	ResendVerificationEmail(ctx *fiber.Ctx, email string) error
}

// EmailVerification is the service that handles email verification logic.
type EmailVerification struct {
	VerificationRepo repositories.IEmailVerification
	UserRepo         repositories.IUsersRepository
}

// NewEmailVerificationService creates a new instance of the EmailVerification service.
func NewEmailVerificationService(verificationRepository repositories.IEmailVerification, userRepository repositories.IUsersRepository) IEmailVerificationService {
	if verificationRepository == nil {
		log.Fatal("nil verification repository")
	}
	if userRepository == nil {
		log.Fatal("nil user repository")
	}
	return &EmailVerification{
		VerificationRepo: verificationRepository,
		UserRepo:         userRepository,
	}
}

// CreateVerificationRequest handles the logic for creating and sending an email verification request.
func (s *EmailVerification) CreateVerificationRequest(ctx *fiber.Ctx, email string) error {
	fmt.Println("Attempting to send verification email to: ", email)

	// Fetch user information based on the provided email
	user, err := s.UserRepo.GetUserInfo(email)
	fmt.Println("got user info")
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user of this email does not exist")
	}

	// Generate a secure token for the verification link
	token, err := GenerateToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create the entity to store in the database
	verificationReq := &entities.EmailVerificationModel{
		Id:          uuid.NewString(),
		UserId:      user.UserID.String(),
		Email:       user.Email,
		Token:       token,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		IsVerify:    false,
		IsAvailable: true,
	}

	// Save the verification request to the database
	err = s.VerificationRepo.InsertNewEmailVerify(context.Background(), verificationReq)
	if err != nil {
    
    log.Printf("[EmailVerification] insert failed: %+v", err)
    return fmt.Errorf("failed to save verification request: %w", err)
}
	fmt.Println("inserted")
	// Construct the verification link
	verificationLink := fmt.Sprintf(os.Getenv("FRONTEND_URL")+"/verifyemail?token=%s", token)

	// Create the simple HTML email body
	emailBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body>
		<div style="font-family: Arial, sans-serif;">
			<h2>Email Verification</h2>
			<p>Hello %s,</p>
			<p>Thank you for registering. Please click the button below to verify your email address:</p>
			<p>
				<a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Verify Email</a>
			</p>
			<p>This link will expire in 24 hours.</p>
			<br>
			<p style="font-size: 12px; color: #888;">If you did not sign up, you can ignore this email.</p>
		</div>
	</body>
	</html>
	`, user.Email, verificationLink)

	// Send the email
	err = SendEmail(
		"smtp.gmail.com",
		587,
		email,
		"Verify Your Email Address",
		emailBody,
	)
	if err != nil {
		fmt.Println("failed to send verification email")
	}
	fmt.Println("email sent")

	return nil
}

func (s *EmailVerification) VerifyEmail(token string) error {
	if token == "" {
		return fmt.Errorf("no token")
	}
	data, err := s.VerificationRepo.GetEmailVerificationByToken(context.Background(), token)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error get reset data")
	}
	if data == nil {
		return fmt.Errorf("no token match")
	}
	if !data.IsAvailable || time.Now().After(data.ExpiresAt) {
		return fmt.Errorf("this token is unavailable")
	}

	data.IsAvailable = false
	data.IsVerify = true

	err = s.VerificationRepo.UpdateEmailVerification(context.Background(), data)
	if err != nil {
		fmt.Println("Error to update")
		return fmt.Errorf("error to update")
	}
	if data.UserId != "" {
        err = s.UserRepo.SetUserVerify(data.UserId, true)
        if err != nil {
            fmt.Println("Error to update user verified status:", err)
            return fmt.Errorf("error to update user verified status")
        }
    }

	return nil
}

func (s *EmailVerification) ResendVerificationEmail(ctx *fiber.Ctx, email string) error {
	fmt.Println("resending email")
	if s.UserRepo == nil || s.VerificationRepo == nil {
		return fmt.Errorf("email verification service not wired correctly (nil repo)")
	}

	if email == "" {
		fmt.Println("email not found")
		return fmt.Errorf("email not found")
	}
	user, err := s.UserRepo.GetUserInfo(email)
	fmt.Print("im here")
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user of this email does not exist")
	}
	if err := s.VerificationRepo.DeactivateAllByEmail(context.Background(), email); err != nil {
		return fmt.Errorf("failed to deactivate previous requests: %w", err)
	}
	return s.CreateVerificationRequest(ctx, email)
}
