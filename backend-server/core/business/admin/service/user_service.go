package service

import "github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"

type UserService interface {
	GetAllUsers(page int64, pageSize int64, sort string, filter *user.UserFilterDto) ([]*user.UserDto, int64, error)
	CreateUser(userDto *user.UserCreateDto) error
	GetUserByID(userID int64) (*user.UserDto, error)
	UpdateUser(userID int64, user *user.UserUpdateDto) error
	DeleteUser(userID int64) error
	DeleteUserByIDs(userIDs []int64) error

	GetByUsername(username string) (*user.UserDto, error)
	GetByEmail(email string) (*user.UserDto, error)
	GetByStaffID(staffID string) (*user.UserDto, error)
}
