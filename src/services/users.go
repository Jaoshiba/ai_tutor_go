package services

import (
	"errors"
	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"
	util "go-fiber-template/src/services/utils"
	"time"
	"github.com/google/uuid"
	"fmt"
)

type usersService struct {
	UsersRepository repositories.IUsersRepository // <-- PostgreSQL implementation
}

type IUsersService interface {
	GetAllUsers() (*[]entities.UserDataModel, error)
	InsertNewUser(data entities.UserDataModel) error
	CheckByEmail(email string) error
	GetUserPublicProfile(userId string) (*entities.UserInfoModel, error)
}

func NewUsersService(repo0 repositories.IUsersRepository) IUsersService {
	return &usersService{
		UsersRepository: repo0,
	}
}

func (sv *usersService) GetAllUsers() (*[]entities.UserDataModel, error) {
	return sv.UsersRepository.FindAll()
}

func (sv *usersService) InsertNewUser(data entities.UserDataModel) error {

	if err := sv.CheckByEmail(data.Email); err != nil {
		return err // จะหยุดไม่ insert ถ้า email ซ้ำ
	}
	if err := sv.CheckByUsername(data.Username); err != nil {
		return err // จะหยุดไม่ insert ถ้า email ซ้ำ
	}
	data.CreatedAt = time.Now().Add(7 * time.Hour)
	data.UpdatedAt = time.Now().Add(7 * time.Hour)

	data.UserID = uuid.New()

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return err
	}
	data.Password = hashedPassword

	return sv.UsersRepository.InsertUser(data)
}

func (sv *usersService) CheckByEmail(email string) error {
	user, err := sv.UsersRepository.FindByEmail(email)
	if err != nil {
		return err // error จาก database
	}
	if user != nil {
		return errors.New("email already exists")
	}
	return nil // email ใช้งานได้
}

func (sv *usersService) CheckByUsername(username string) error {
	user, err := sv.UsersRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("email already exists")
	}
	return nil // email ใช้งานได้
}

func (sv *usersService) CheckUserExistBy(field string, value string) error {
	var user *entities.UserDataModel
	var err error

	switch field {
	case "email":
		user, err = sv.UsersRepository.FindByEmail(value)
	case "username":
		user, err = sv.UsersRepository.FindByUsername(value)
	default:
		return errors.New("invalid field for user lookup")
	}

	if err != nil {
		return err
	}

	if user != nil {
		return fmt.Errorf("%s already exists", field)
	}

	return nil
}
func (sv *usersService) UpdateUserProfile(req entities.UpdateUserProfileRequest) error {
	return sv.UsersRepository.UpdateUserProfile(req)
}

func (sv *usersService) GetUserPublicProfile(userId string) (*entities.UserInfoModel, error) {
    user, err := sv.UsersRepository.GetUserById(userId)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    // map UserDataModel -> UserInfoModel
    return &entities.UserInfoModel{
        UserID:          user.UserID,
        Username:        user.Username,
        FirstName:       user.FirstName,
        LastName:        user.LastName,
        Email:           user.Email,
        Gender:          user.Gender,
        Role:            user.Role,
        DOB:             user.DOB,
        CreatedAt:       user.CreatedAt,
        UpdatedAt:       user.UpdatedAt,
        IsEmailVerified: user.IsEmailVerified,
    }, nil
}
