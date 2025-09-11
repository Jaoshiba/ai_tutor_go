package repositories

import (
	"context"
	"database/sql"
	"log"
	"fmt"

	"go-fiber-template/domain/entities"
)

type IResetPassword interface {
	InsertNewResetPassword(ctx context.Context, data *entities.ResetPassword) error
	GetResetPasswordByToken(ctx context.Context, token string) (*entities.ResetPassword, error)
	UpdateResetPasswordStatus(ctx context.Context, id string, isReset bool) error
	DeleteResetPassword(ctx context.Context, id string) error
}

type resetPasswordRepository struct {
	db *sql.DB
}

func NewResetPasswordRepository(db *sql.DB) IResetPassword {
	if db == nil {
		log.Fatal("got nil db")
	}
	return &resetPasswordRepository{
		db: db,
	}
}

func (repo *resetPasswordRepository) InsertNewResetPassword(ctx context.Context, data *entities.ResetPassword) error {
	query := `
		INSERT INTO reset_passwords (id, user_id, created_at, expires_at, token, is_reset, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := repo.db.ExecContext(ctx, query,
		data.Id,
		data.UserId,
		data.CreatedAt,
		data.ExpiresAt,
		data.Token,
		data.IsReset,
	)

	if err != nil {
		return fmt.Errorf("failed to insert new reset password record: %w", err)
	}
	return nil
}

func (repo *resetPasswordRepository) GetResetPasswordByToken(ctx context.Context, token string) (*entities.ResetPassword, error) {
	query := `
		SELECT id, user_id, created_at, expires_at, token, is_reset, is_available
		FROM reset_passwords
		WHERE token = $1
	`
	row := repo.db.QueryRowContext(ctx, query, token)

	data := &entities.ResetPassword{}
	err := row.Scan(
		&data.Id,
		&data.UserId,
		&data.CreatedAt,
		&data.ExpiresAt,
		&data.Token,
		&data.IsReset,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reset password record not found with token %s: %w", token, err)
		}
		return nil, fmt.Errorf("failed to get reset password record by token: %w", err)
	}
	return data, nil
}

func (repo *resetPasswordRepository) UpdateResetPasswordStatus(ctx context.Context, id string, isReset bool) error {
	query := `
		UPDATE reset_passwords
		SET is_reset = $1,
			is_available = false
		WHERE id = $2
	`
	_, err := repo.db.ExecContext(ctx, query, isReset, id)
	if err != nil {
		return fmt.Errorf("failed to update reset password status: %w", err)
	}
	return nil
}

func (repo *resetPasswordRepository) DeleteResetPassword(ctx context.Context, id string) error {
	query := `
		DELETE FROM reset_passwords
		WHERE id = $1
	`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reset password record: %w", err)
	}
	return nil
}