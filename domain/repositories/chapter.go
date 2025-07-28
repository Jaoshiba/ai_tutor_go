package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
)

type chaptersRepository struct {
	db *sql.DB
}

type IChapterRepository interface {
	InsertChapter(chapter entities.ChapterDataModel) error
}

func NewChapterRepository(db *sql.DB) IChapterRepository {
	return &chaptersRepository{
		db: db,
	}
}

func (repo *chaptersRepository) InsertChapter(chapter entities.ChapterDataModel) error {

	fmt.Println("InsertChapter called with chapter:", chapter)
	query := `
	INSERT INTO chapters (
		chapterid, chaptername, userid, courseid, chaptercontent, isfinished, createat, updateat
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	fmt.Println("userid in chap repo : ", chapter.UserID)
	_, err := repo.db.ExecContext(context.Background(), query,
		chapter.ChapterId,
		chapter.ChapterName,
		chapter.UserID,
		chapter.CouseId,
		chapter.ChapterContent,
		chapter.IsFinished,
		chapter.CreateAt,
		chapter.UpdatedAt,
	)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
