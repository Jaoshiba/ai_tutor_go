package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go-fiber-template/domain/entities"
)

type emailVerificationRepository struct {
	db *sql.DB
}

type IEmailVerification interface {
	InsertNewEmailVerify(ctx context.Context, data *entities.EmailVerificationModel) error
	GetEmailVerificationById(ctx context.Context, id string) (*entities.EmailVerificationModel, error)
	GetEmailVerificationByToken(ctx context.Context, token string) (*entities.EmailVerificationModel, error)
	UpdateEmailVerification(ctx context.Context, data *entities.EmailVerificationModel) error
	DeleteEmailVerification(ctx context.Context, id string) error
	DeactivateAllByEmail(ctx context.Context, email string) error
}

func NewEmailVerificationRepository(db *sql.DB) IEmailVerification {
	if db == nil {
		log.Fatal("got nil db")
	}
	fmt.Println("received DB", db)
	return &emailVerificationRepository{
		db: db,
	}
}

func (repo *emailVerificationRepository) InsertNewEmailVerify(ctx context.Context, data *entities.EmailVerificationModel) error {
	query := `
		INSERT INTO email_verifications (id, user_id, email, created_at, expires_at, token, is_verify, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := repo.db.ExecContext(ctx, query,
		data.Id,
		data.UserId,
		data.Email,
		data.CreatedAt,
		data.ExpiresAt,
		data.Token,
		data.IsVerify,
		data.IsAvailable,
	)

	if err != nil {
		return fmt.Errorf("failed to insert new email verification: %w", err)
	}
	return nil
}

func (repo *emailVerificationRepository) GetEmailVerificationById(ctx context.Context, id string) (*entities.EmailVerificationModel, error) {
	query := `
		SELECT id, user_id, email, created_at, expires_at, token, is_verify, is_available
		FROM email_verifications
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)

	data := &entities.EmailVerificationModel{}
	err := row.Scan(
		&data.Id,
		&data.UserId,
		&data.Email,
		&data.CreatedAt,
		&data.ExpiresAt,
		&data.Token,
		&data.IsVerify,
		&data.IsAvailable,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email verification not found with id %s: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get email verification by id: %w", err)
	}
	return data, nil
}

func (repo *emailVerificationRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*entities.EmailVerificationModel, error) {
	query := `
		SELECT id, user_id, email, created_at, expires_at, token, is_verify, is_available
		FROM email_verifications
		WHERE token = $1
	`
	row := repo.db.QueryRowContext(ctx, query, token)

	data := &entities.EmailVerificationModel{}
	err := row.Scan(
		&data.Id,
		&data.UserId,
		&data.Email,
		&data.CreatedAt,
		&data.ExpiresAt,
		&data.Token,
		&data.IsVerify,
		&data.IsAvailable,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email verification not found with token %s: %w", token, err)
		}
		return nil, fmt.Errorf("failed to get email verification by token: %w", err)
	}
	return data, nil
}

func (repo *emailVerificationRepository) UpdateEmailVerification(ctx context.Context, data *entities.EmailVerificationModel) error {
	query := `
		UPDATE email_verifications
SET is_verify = $1,
    updated_at = NOW(),
    is_available = $3
WHERE id = $2

	`
	_, err := repo.db.ExecContext(ctx, query, data.IsVerify, data.Id, data.IsAvailable)
	if err != nil {
		return fmt.Errorf("failed to update email verification: %w", err)
	}
	return nil
}

func (repo *emailVerificationRepository) DeleteEmailVerification(ctx context.Context, id string) error {
	query := `
		DELETE FROM email_verifications
		WHERE id = $1
	`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete email verification: %w", err)
	}
	return nil
}
func (repo *emailVerificationRepository) DeactivateAllByEmail(ctx context.Context, email string) error {
	query := `
		UPDATE email_verifications
		SET is_available = FALSE,
		    updated_at = NOW()
		WHERE email = $1
		  AND is_available = TRUE
		  AND is_verify = FALSE
	`
	_, err := repo.db.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to deactivate old verification requests: %w", err)
	}
	return nil
}
