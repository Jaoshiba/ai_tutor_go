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
			chapterid, chaptername, userid, roadmapid, chaptercontent, isfinished, createdat, updatedat
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := repo.db.ExecContext(context.Background(), query,
		chapter.ChapterId,
		chapter.ChapterName,
		chapter.UserID,
		chapter.RoadmapId,
		chapter.ChapterContent,
		chapter.IsFinished,
		chapter.CreateAt,
		chapter.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
