package service

import (
	"github.com/vamika-digital/wms-api-server/internal/auth/dto"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
)

type UserService interface {
	CreateUser(user dto.UserCreateDto) error
	UpdateUser(user dto.UserUpdateDto) error
	DeleteUser(userID int64) error
	GetUserByID(userID int64) (dto.UserDto, error)
	GetAllUsers(page int, pageSize int, sort string, filter repository.UserFilterOptions) ([]dto.UserDto, int, error)
	GetUserByUsername(username string) (dto.UserDto, error)
	GetUserByEmail(email string) (dto.UserDto, error)
}
