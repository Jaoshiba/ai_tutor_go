package repositories

import (
	"context"
	"database/sql"
	"go-fiber-template/domain/entities"
	"log"
)

// PostgreSQL implementation
type usersRepositoryPostgres struct {
	db *sql.DB
}

type IUsersRepository interface {
	InsertUser(data entities.UserDataModel) error
	FindAll() (*[]entities.UserDataModel, error)
	FindByEmail(email string) (*entities.UserDataModel, error)
	FindByUsername(username string) (*entities.UserDataModel, error)
}

func NewUsersRepositoryPostgres(db *sql.DB) IUsersRepository {
	return &usersRepositoryPostgres{
		db: db,
	}
}

func (repo *usersRepositoryPostgres) InsertUser(data entities.UserDataModel) error {
	query := `
		INSERT INTO users (
			userid, username, firstname, lastname,
			email, gender, role, dob,
			password, createdat, updateat
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := repo.db.ExecContext(context.Background(), query,
		data.UserID,
		data.Username,
		data.FirstName,
		data.LastName,
		data.Email,
		data.Gender,
		data.Role,
		data.DOB,
		data.Password,
		data.CreatedAt,
		data.UpdatedAt,
	)

	if err != nil {
		log.Printf("Users -> InsertUser: %v\n", err)
		return err
	}
	return nil
}


func (repo *usersRepositoryPostgres) FindAll() (*[]entities.UserDataModel, error) {
	query := `SELECT id, name, email, role, picture FROM users`
	rows, err := repo.db.QueryContext(context.Background(), query)
	if err != nil {
		log.Printf("Users -> FindAll: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []entities.UserDataModel
	for rows.Next() {
		var user entities.UserDataModel
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			log.Printf("Users -> Scan: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (repo *usersRepositoryPostgres) FindByEmail(email string) (*entities.UserDataModel, error) {
	query := `SELECT userid, username, email, password, role FROM users WHERE email = $1 LIMIT 1`
	row := repo.db.QueryRowContext(context.Background(), query, email)

	var user entities.UserDataModel
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่พบผู้ใช้ ถือว่า OK
		}
		return nil, err // เกิด error จริง
	}
	return &user, nil
}


func (repo *usersRepositoryPostgres) FindByUsername(username string) (*entities.UserDataModel, error) {
	query := `SELECT userid, username, email, password, role FROM users WHERE username = $1 LIMIT 1`
	row := repo.db.QueryRowContext(context.Background(), query, username)

	var user entities.UserDataModel
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
