// repositories/sessions/session_repository.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"

    "go-fiber-template/domain/entities"
)

type ISessionRepository interface {
    CreateSession(session entities.SessionDataModel) error
    GetSessionByID(sessionID string) (*entities.SessionDataModel, error)
    UpdateSession(session entities.SessionDataModel) error
    DeleteSession(sessionID string) error
    ListSessionsByUser(userID string) ([]entities.SessionDataModel, error)
}

type sessionRepository struct {
    db *sql.DB
}

func NewSessionRepository(db *sql.DB) ISessionRepository {
    return &sessionRepository{db: db}
}

// CREATE
func (r *sessionRepository) CreateSession(session entities.SessionDataModel) error {
    query := `
        INSERT INTO sessions (sessionid, userid, createdat, expiredat)
        VALUES ($1, $2, $3, $4)
    `
    _, err := r.db.ExecContext(context.Background(), query,
        session.SessionID,
        session.UserID,
        session.CreatedAt,
        session.ExpiredAt,
    )
    if err != nil {
        return fmt.Errorf("CreateSession error: %w", err)
    }
    return nil
}

// READ
func (r *sessionRepository) GetSessionByID(sessionID string) (*entities.SessionDataModel, error) {
    query := `
        SELECT sessionid, userid, createdat, expiredat
        FROM sessions
        WHERE sessionid=$1
    `
    s := &entities.SessionDataModel{}
    err := r.db.QueryRowContext(context.Background(), query, sessionID).Scan(
        &s.SessionID,
        &s.UserID,
        &s.CreatedAt,
        &s.ExpiredAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("GetSessionByID error: %w", err)
    }
    return s, nil
}

// UPDATE
func (r *sessionRepository) UpdateSession(session entities.SessionDataModel) error {
    query := `
        UPDATE sessions
        SET userid=$1, createdat=$2, expiredat=$3
        WHERE sessionid=$4
    `
    _, err := r.db.ExecContext(context.Background(), query,
        session.UserID,
        session.CreatedAt,
        session.ExpiredAt,
        session.SessionID,
    )
    if err != nil {
        return fmt.Errorf("UpdateSession error: %w", err)
    }
    return nil
}

// DELETE
func (r *sessionRepository) DeleteSession(sessionID string) error {
    query := `
        DELETE FROM sessions
        WHERE sessionid=$1
    `
    _, err := r.db.ExecContext(context.Background(), query, sessionID)
    if err != nil {
        return fmt.Errorf("DeleteSession error: %w", err)
    }
    return nil
}

// LIST SESSIONS BY USER
func (r *sessionRepository) ListSessionsByUser(userID string) ([]entities.SessionDataModel, error) {
    query := `
        SELECT sessionid, userid, createdat, expiredat
        FROM sessions
        WHERE userid=$1
        ORDER BY createdat DESC
    `
    rows, err := r.db.QueryContext(context.Background(), query, userID)
    if err != nil {
        return nil, fmt.Errorf("ListSessionsByUser error: %w", err)
    }
    defer rows.Close()

    sessions := []entities.SessionDataModel{}
    for rows.Next() {
        s := entities.SessionDataModel{}
        if err := rows.Scan(&s.SessionID, &s.UserID, &s.CreatedAt, &s.ExpiredAt); err != nil {
            return nil, fmt.Errorf("ListSessionsByUser scan error: %w", err)
        }
        sessions = append(sessions, s)
    }
    return sessions, nil
}
