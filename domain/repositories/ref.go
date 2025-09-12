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
	InsertRef(ref entities.SearchLinks) error
	GetRefsByModuleId(moduleId string) ([]entities.SearchLinks, error)
}

func NewRefRepository(db *sql.DB) IRefInterface {
	return &refRepository{
		db: db,
	}
}

func (rr *refRepository) InsertRef(ref entities.SearchLinks) error {

	fmt.Println("InsertRef called with ref:", ref)

	query := `
		INSERT INTO moduleref (
			linkid, moduleid, title, link, snipppet, searchat
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := rr.db.ExecContext(context.Background(), query,
		ref.LinkID,
		ref.ModuleID,
		ref.Title,
		ref.Link,
		ref.Snippet,
		ref.Searchat,
	)
	if err != nil {
		return err
	}
	return nil
}

func (rr *refRepository) GetRefsByModuleId(moduleId string) ([]entities.SearchLinks, error) {
	query := `
		SELECT refid, moduleid, title, link, content, searchat
		FROM moduleref
		WHERE moduleid = $1`

	rows, err := rr.db.QueryContext(context.Background(), query, moduleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []entities.SearchLinks
	for rows.Next() {
		var ref entities.SearchLinks
		if err := rows.Scan(
			&ref.LinkID,
			&ref.ModuleID,
			&ref.Title,
			&ref.Link,
			&ref.Snippet,
			&ref.Searchat,
		); err != nil {
			return nil, err
		}
		refs = append(refs, ref)

	}

	return refs, nil

}
