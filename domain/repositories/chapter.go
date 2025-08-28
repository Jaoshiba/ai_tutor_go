package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
	"log"
)

type chaptersRepository struct {
	db *sql.DB
}

type IChapterRepository interface {
	InsertChapter(chapter entities.ChapterDataModel) error
	GetChaptersByModuleID(moduleID string) ([]entities.ChapterDataModel, error)
	DeleteChapter(chapterID string) error
	DeleteChapterByModuleID(moduleID string) error

}

func NewChapterRepository(db *sql.DB) IChapterRepository {
	if db == nil {
		log.Fatal("❌ NewChapterRepository got nil DB")
	}
	fmt.Println("✅ NewChapterRepository received DB:", db)
	return &chaptersRepository{
		db: db,
	}
}

func (repo *chaptersRepository) InsertChapter(chapter entities.ChapterDataModel) error {

	fmt.Println("InsertChapter called with chapter:", chapter)
	query := `
	INSERT INTO chapters (
		chapterid, chaptername, userid, courseid, chaptercontent, isfinished, createat, updateat, moduleid
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	fmt.Println("userid in chap repo : ", chapter.UserID)
	_, err := repo.db.ExecContext(context.Background(), query,
		chapter.ChapterId,
		chapter.ChapterName,
		chapter.UserID,
		chapter.CourseId,
		chapter.ChapterContent,
		chapter.IsFinished,
		chapter.CreateAt,
		chapter.UpdatedAt,
		chapter.ModuleId,
	)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func (repo *chaptersRepository) GetChaptersByModuleID(moduleID string) ([]entities.ChapterDataModel, error) {
	fmt.Println("im in repo")
	query := `
		SELECT chapterid, chaptername, userid, courseid, chaptercontent, isfinished, createat, updateat, moduleid
		FROM chapters
		WHERE moduleid = $1
		ORDER BY createat
	`
	rows, err := repo.db.QueryContext(context.Background(), query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chapters by module ID %s: %w", moduleID, err)
	}
	defer rows.Close()

	var chapters []entities.ChapterDataModel
	for rows.Next() {
		var chapter entities.ChapterDataModel
		// สแกนตรงเข้าสู่ฟิลด์ time.Time ได้เลย
		if err := rows.Scan(
			&chapter.ChapterId,
			&chapter.ChapterName,
			&chapter.UserID,
			&chapter.CourseId,
			&chapter.ChapterContent,
			&chapter.IsFinished,
			&chapter.CreateAt,  // สแกนตรงนี้
			&chapter.UpdatedAt, // สแกนตรงนี้
			&chapter.ModuleId,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chapter row: %w", err)
		}

		// ไม่ต้องมีการแปลง string เป็น time.Time อีกแล้ว
		chapters = append(chapters, chapter)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during chapters row iteration: %w", err)
	}

	return chapters, nil
}

// chapters_repo.go
func (repo *chaptersRepository) DeleteChapterByModuleID(moduleID string) error {
    if moduleID == "" {
        return fmt.Errorf("moduleID is empty")
    }
    const q = `DELETE FROM chapters WHERE moduleid = $1`
    result, err := repo.db.ExecContext(context.Background(), q, moduleID)
    if err != nil {
        return fmt.Errorf("delete chapters by module_id failed: %w", err)
    }
    // ไม่ต้อง error เมื่อ 0 แถว
    if n, _ := result.RowsAffected(); n > 0 {
        fmt.Printf("Deleted %d chapter(s) for module_id=%s\n", n, moduleID)
    }
    return nil
}


func (repo *chaptersRepository) DeleteChapter(chapterID string) error {
    fmt.Printf("DeleteChapter called for chapter with ID: %s\n", chapterID)
    
    query := `
        DELETE FROM chapters
        WHERE chapterid = $1
    `

    result, err := repo.db.ExecContext(context.Background(), query, chapterID)
    if err != nil {
        return fmt.Errorf("failed to delete chapter with ID %s: %w", chapterID, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected after deleting chapter: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no chapter found with ID %s to delete", chapterID)
    }
    
    fmt.Printf("Successfully deleted %d row for chapter ID: %s\n", rowsAffected, chapterID)
    return nil
}