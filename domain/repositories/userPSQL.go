package repositories

import (
	"context"
	"database/sql"

	// "fmt"
	"go-fiber-template/domain/entities"
	"log"

	"github.com/google/uuid"
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
	GetUserInfo(email string) (*entities.UserInfoModel,error)
	UpdateUserProfile(data entities.UpdateUserProfileRequest) error
	UpdateUserPassword(userID string, newHashedPassword string) error

	GetUserById(userId string) (*entities.UserDataModel, error)
	SetUserVerify(userID string, verified bool) error
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
			password, createdat, updatedat
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
	query := `SELECT id, name, email, role, isemailverified FROM users`
	rows, err := repo.db.QueryContext(context.Background(), query)
	if err != nil {
		log.Printf("Users -> FindAll: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []entities.UserDataModel
	for rows.Next() {
		var user entities.UserDataModel
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Role, &user.IsEmailVerified)
		if err != nil {
			log.Printf("Users -> Scan: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (repo *usersRepositoryPostgres) FindByEmail(email string) (*entities.UserDataModel, error) {
	query := `SELECT userid, username, email, password, role, isemailverified is FROM users WHERE email = $1 LIMIT 1`
	row := repo.db.QueryRowContext(context.Background(), query, email)

	var user entities.UserDataModel
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role, &user.IsEmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่พบผู้ใช้ ถือว่า OK
		}
		return nil, err // เกิด error จริง
	}
	return &user, nil
}

func (repo *usersRepositoryPostgres) GetUserInfo(email string) (*entities.UserInfoModel,error) {
	query := `SELECT userid, username, email, gender, isemailverified, firstname, lastname, role, dob, createdat FROM users WHERE email = $1 LIMIT 1`
	row := repo.db.QueryRowContext(context.Background(), query, email)
	
	var user entities.UserInfoModel
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Gender, &user.IsEmailVerified, &user.FirstName, &user.LastName, &user.Role, &user.DOB, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err 
	}
	return &user, nil
}

func (repo *usersRepositoryPostgres) GetUserById(userId string) (*entities.UserDataModel, error) {
    query := `SELECT userid, username, firstname, lastname, gender, role, dob, isemailverified, createdat FROM users WHERE userid = $1 LIMIT 1`
    row := repo.db.QueryRowContext(context.Background(), query, userId)
    var user entities.UserDataModel
    err := row.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.Role, &user.DOB, &user.IsEmailVerified, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}



func (repo *usersRepositoryPostgres) FindByUsername(username string) (*entities.UserDataModel, error) {
	query := `SELECT userid, username, email, password, role, isemailverified FROM users WHERE username = $1 LIMIT 1`
	row := repo.db.QueryRowContext(context.Background(), query, username)

	var user entities.UserDataModel
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Role, &user.IsEmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *usersRepositoryPostgres) UpdateUserProfile(data entities.UpdateUserProfileRequest) error {
	// data.UserID เป็น uuid.UUID ตาม entities
	// data.DOB เป็น string: ให้แคสต์เป็น date ถ้าไม่ว่าง, ถ้าว่างจะเป็น NULL
	query := `
		UPDATE users
		SET firstname = $1,
		    lastname  = $2,
		    gender    = $3,
		    dob       = NULLIF($4, '')::date,
		    updatedat = NOW()
		WHERE userid = $5
	`

	// ป้องกัน panic ถ้าฝั่ง call site เผลอส่ง zero UUID มา
	var userID uuid.UUID = data.UserID

	res, err := repo.db.ExecContext(
		context.Background(),
		query,
		data.FirstName,
		data.LastName,
		data.Gender,
		data.DOB,  
		userID,    
	)
	if err != nil {
		log.Printf("Users -> UpdateUserProfile: %v\n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Users -> UpdateUserProfile RowsAffected error: %v\n", err)
		// เลือกรีเทิร์นต่อไป เพราะอัปเดตน่าจะสำเร็จแล้ว แค่นับแถวไม่ได้
	}

	if rows == 0 {
		// ไม่พบผู้ใช้ตาม userid
		return sql.ErrNoRows
	}
	return nil
}


func (repo *usersRepositoryPostgres) UpdateUserPassword(userID string, newHashedPassword string) error {
	query := `
		UPDATE users
		SET password  = $1,
		    updatedat = NOW()
		WHERE userid = $2
	`

	res, err := repo.db.ExecContext(
		context.Background(),
		query,
		newHashedPassword,
		userID,
	)
	if err != nil {
		log.Printf("Users -> UpdateUserPassword: %v\n", err)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *usersRepositoryPostgres) SetUserVerify(userID string, verified bool) error {
	query := `
		UPDATE users
		SET isemailverified = $1,
		    updatedat       = NOW()
		WHERE userid = $2
	`
	res, err := repo.db.ExecContext(context.Background(), query, verified, userID)
	if err != nil {
		log.Printf("Users -> SetUserVerify: %v\n", err)
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *usersRepositoryPostgres) UpdateUserProfileExceptEmail(data entities.UpdateUserProfileRequest) error {
    query := `
        UPDATE users
        SET firstname = $1,
            lastname  = $2,
            gender    = $3,
            dob       = $4,
            updatedat = NOW()
        WHERE userid = $5
    `
    res, err := repo.db.ExecContext(
        context.Background(),
        query,
        data.FirstName,
        data.LastName,
        data.Gender,
        data.DOB,
        data.UserID,
    )
    if err != nil {
        log.Printf("Users -> UpdateUserProfileExceptEmail: %v\n", err)
        return err
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return sql.ErrNoRows
    }
    return nil
}
