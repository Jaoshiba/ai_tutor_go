// repositories/learningprogress/learningprogress_repository.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"

    "go-fiber-template/domain/entities"
)

type ILearningProgressRepository interface {
    CreateProgress(progress entities.LearningProgressDataModel) error
    GetProgressByID(processID string) (*entities.LearningProgressDataModel, error)
    UpdateProgress(progress entities.LearningProgressDataModel) error
    DeleteProgress(processID string) error
    ListProgressByUser(userID string) ([]entities.LearningProgressDataModel, error)
    ListProgressByModule(moduleID string) ([]entities.LearningProgressDataModel, error)
}

type learningProgressRepository struct {
    db *sql.DB
}

func NewLearningProgressRepository(db *sql.DB) ILearningProgressRepository {
    return &learningProgressRepository{db: db}
}

// CREATE
func (r *learningProgressRepository) CreateProgress(progress entities.LearningProgressDataModel) error {
    query := `
        INSERT INTO learningprogress (processid, userid, moduleid, chapterid, createdat)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := r.db.ExecContext(context.Background(), query,
        progress.ProcessID,
        progress.UserID,
        progress.ModuleID,
        progress.ChapterID,
        progress.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("CreateProgress error: %w", err)
    }
    return nil
}

// READ
func (r *learningProgressRepository) GetProgressByID(processID string) (*entities.LearningProgressDataModel, error) {
    query := `
        SELECT processid, userid, moduleid, chapterid, createdat
        FROM learningprogress
        WHERE processid=$1
    `
    p := &entities.LearningProgressDataModel{}
    err := r.db.QueryRowContext(context.Background(), query, processID).Scan(
        &p.ProcessID,
        &p.UserID,
        &p.ModuleID,
        &p.ChapterID,
        &p.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("GetProgressByID error: %w", err)
    }
    return p, nil
}

// UPDATE
func (r *learningProgressRepository) UpdateProgress(progress entities.LearningProgressDataModel) error {
    query := `
        UPDATE learningprogress
        SET userid=$1, moduleid=$2, chapterid=$3, createdat=$4
        WHERE processid=$5
    `
    _, err := r.db.ExecContext(context.Background(), query,
        progress.UserID,
        progress.ModuleID,
        progress.ChapterID,
        progress.CreatedAt,
        progress.ProcessID,
    )
    if err != nil {
        return fmt.Errorf("UpdateProgress error: %w", err)
    }
    return nil
}

// DELETE
func (r *learningProgressRepository) DeleteProgress(processID string) error {
    query := `
        DELETE FROM learningprogress
        WHERE processid=$1
    `
    _, err := r.db.ExecContext(context.Background(), query, processID)
    if err != nil {
        return fmt.Errorf("DeleteProgress error: %w", err)
    }
    return nil
}

// LIST BY USER
func (r *learningProgressRepository) ListProgressByUser(userID string) ([]entities.LearningProgressDataModel, error) {
    query := `
        SELECT processid, userid, moduleid, chapterid, createdat
        FROM learningprogress
        WHERE userid=$1
        ORDER BY createdat DESC
    `
    rows, err := r.db.QueryContext(context.Background(), query, userID)
    if err != nil {
        return nil, fmt.Errorf("ListProgressByUser error: %w", err)
    }
    defer rows.Close()

    progresses := []entities.LearningProgressDataModel{}
    for rows.Next() {
        p := entities.LearningProgressDataModel{}
        if err := rows.Scan(&p.ProcessID, &p.UserID, &p.ModuleID, &p.ChapterID, &p.CreatedAt); err != nil {
            return nil, fmt.Errorf("ListProgressByUser scan error: %w", err)
        }
        progresses = append(progresses, p)
    }
    return progresses, nil
}

// LIST BY MODULE
func (r *learningProgressRepository) ListProgressByModule(moduleID string) ([]entities.LearningProgressDataModel, error) {
    query := `
        SELECT processid, userid, moduleid, chapterid, createdat
        FROM learningprogress
        WHERE moduleid=$1
        ORDER BY createdat DESC
    `
    rows, err := r.db.QueryContext(context.Background(), query, moduleID)
    if err != nil {
        return nil, fmt.Errorf("ListProgressByModule error: %w", err)
    }
    defer rows.Close()

    progresses := []entities.LearningProgressDataModel{}
    for rows.Next() {
        p := entities.LearningProgressDataModel{}
        if err := rows.Scan(&p.ProcessID, &p.UserID, &p.ModuleID, &p.ChapterID, &p.CreatedAt); err != nil {
            return nil, fmt.Errorf("ListProgressByModule scan error: %w", err)
        }
        progresses = append(progresses, p)
    }
    return progresses, nil
}
