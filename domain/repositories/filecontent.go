package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "go-fiber-template/domain/entities"
)

type filecontentRepository struct {
    db *sql.DB
}

type IFileContentRepository interface {
    InsertFileContent(ctx context.Context, content entities.FileContentModel) error
}

func NewFileContentRepository(db *sql.DB) IFileContentRepository {
    return &filecontentRepository{
        db: db,
    }
}

func (r *filecontentRepository) InsertFileContent(ctx context.Context, content entities.FileContentModel) error {
    fmt.Println("InsertFileContent called with file content:", content)

    query := `
        INSERT INTO file_content (
            id, courseid, content, file_path, createat
        ) VALUES ($1, $2, $3, $4, $5)
    `

    // Prepare the SQL statement for execution
    stmt, err := r.db.PrepareContext(ctx, query)
    if err != nil {
        return fmt.Errorf("could not prepare statement: %w", err)
    }
    defer stmt.Close() // Ensure the statement is closed

    // Execute the query with the entity data
    _, err = stmt.ExecContext(ctx, content.Id, content.CourseId, content.Content, content.FilePath, content.CreatedAt)
    if err != nil {
        return fmt.Errorf("could not insert file content: %w", err)
    }

    return nil
}