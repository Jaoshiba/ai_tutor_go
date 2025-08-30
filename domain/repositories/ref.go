package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-template/domain/entities"
)

type refRepository struct {
	db *sql.DB
}

type IRefInterface interface {
	InsertRef(ref entities.RefDataModel) error
}

func NewRefRepository(db *sql.DB) IRefInterface {
	return &refRepository{
		db: db,
	}
}

func (rr *refRepository) InsertRef(ref entities.RefDataModel) error {

	fmt.Println("InsertRef called with ref:", ref)

	query := `
		INSERT INTO moduleref (
			refid, moduleid, title, link, content, searchat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := rr.db.ExecContext(context.Background(), query,
		ref.RefId,
		ref.ModuleId,
		ref.Title,
		ref.Link,
		ref.Content,
		ref.SearchAt,
	)
	if err != nil {
		return err
	}
	return nil
}
