package service

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type AuthService interface {
	GetUserById(idStr string) (*user.UserDto, error)
	ValidateCredentials(username string, password string, loginVia *customtypes.LoginViaEnum) (*user.UserDto, error)
	GenerateAccessToken(user *user.UserDto) (string, error)
	GenerateRefreshToken(user *user.UserDto, expireLong bool) (string, error)
}
