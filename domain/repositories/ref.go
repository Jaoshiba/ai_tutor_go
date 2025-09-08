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
	GetRefsByModuleId(moduleId string) ([]entities.RefDataModel, error)
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

func (rr *refRepository) GetRefsByModuleId(moduleId string) ([]entities.RefDataModel, error) {
	query := `
		SELECT refid, moduleid, title, link, content, searchat
		FROM moduleref
		WHERE moduleid = $1`

	rows, err := rr.db.QueryContext(context.Background(), query, moduleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []entities.RefDataModel
	for rows.Next() {
		var ref entities.RefDataModel
		if err := rows.Scan(
			&ref.RefId,
			&ref.ModuleId,
			&ref.Title,
			&ref.Link,
			&ref.Content,
			&ref.SearchAt,
		); err != nil {
			return nil, err
		}
		refs = append(refs, ref)

	}

	return refs, nil

}
