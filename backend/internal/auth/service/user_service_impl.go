package service

import (
	"github.com/vamika-digital/wms-api-server/internal/auth/domain"
	"github.com/vamika-digital/wms-api-server/internal/auth/dto"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

func (u *UserServiceImpl) CreateUser(userDto dto.UserCreateDto) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := domain.NewUserWithDefaults()
	user.Username = userDto.Username
	user.Password = string(hashedPassword)
	user.Name = userDto.Name
	user.StaffID = userDto.StaffID
	user.Email = userDto.Email
	user.EmailConfirmation = false
	if userDto.EmailConfirmation.Valid {
		user.EmailConfirmation = userDto.EmailConfirmation.Bool
	}
	user.Status = userDto.Status
	user.LastUpdatedBy = userDto.LastUpdatedBy

	user.Roles = userDto.Roles

	userDto.PasswordHash = string(hashedPassword)
	return u.Repo.Create(userDto)
}

func (u *UserServiceImpl) UpdateUser(user dto.UserUpdateDto) error {
	// Check for an existing user with the specified ID
	existingUser, err := u.Repo.GetById(user.ID)
	if err != nil {
		return err
	}

	// Hash the new password if it has been changed
	if user.PasswordHash != "" && user.PasswordHash != existingUser.PasswordHash {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	return u.Repo.Update(user)
}

func (u *UserServiceImpl) DeleteUser(userID int64) error {
	return u.Repo.Delete(userID)
}

func (u *UserServiceImpl) GetUserByID(userID int64) (dto.UserDto, error) {
	return u.Repo.GetById(userID)
}

func (u *UserServiceImpl) GetAllUsers(page int, pageSize int, sort string, filter repository.UserFilterOptions) ([]dto.UserDto, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of users matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UserServiceImpl) GetUserByUsername(username string) (dto.UserDto, error) {
	return u.Repo.FindByUsername(username)
}

func (u *UserServiceImpl) GetUserByEmail(email string) (dto.UserDto, error) {
	return u.Repo.FindByEmail(email)
}
