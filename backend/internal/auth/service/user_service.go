package service

import (
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/user"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
)

type UserService interface {
	CreateUser(user user.UserCreateDto) error
	UpdateUser(user user.UserUpdateDto) error
	DeleteUser(userID int64) error
	GetUserByID(userID int64) (user.UserDto, error)
	GetAllUsers(page int, pageSize int, sort string, filter repository.UserFilterOptions) ([]user.UserDto, int, error)
	GetUserByUsername(username string) (user.UserDto, error)
	GetUserByEmail(email string) (user.UserDto, error)
}
