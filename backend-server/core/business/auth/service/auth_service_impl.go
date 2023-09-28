package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vamika-digital/wms-api-server/config/authconfig"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	adminDomain "github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
	adminRepository "github.com/vamika-digital/wms-api-server/core/business/admin/repository"
	"github.com/vamika-digital/wms-api-server/core/business/auth/dto/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository       adminRepository.UserRepository
	roleRepository       adminRepository.RoleRepository
	permissionRepository adminRepository.PermissionRepository
	userConverter        converter.UserConverter
}

func NewAuthService(userRepository adminRepository.UserRepository, roleRepository adminRepository.RoleRepository, permissionRepository adminRepository.PermissionRepository, userConverter converter.UserConverter) AuthService {
	return &AuthServiceImpl{userRepository: userRepository, roleRepository: roleRepository, permissionRepository: permissionRepository, userConverter: userConverter}
}

func (service *AuthServiceImpl) GetUserById(idStr string) (*user.UserDto, error) {
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, errors.New("invalid datatype")
	}

	var domainUser *adminDomain.User
	domainUser, err = service.userRepository.GetById(userID)

	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	loginUser := service.userConverter.ToDto(domainUser)
	return loginUser, nil
}

func (service *AuthServiceImpl) ValidateCredentials(username string, password string, loginVia *customtypes.LoginViaEnum) (*user.UserDto, error) {
	var domainUser *adminDomain.User
	var err error
	if loginVia == nil || loginVia.ViaEmail() {
		domainUser, err = service.userRepository.GetByEmail(username)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	} else if loginVia == nil || loginVia.ViaStaffID() {
		domainUser, err = service.userRepository.GetByStaffID(username)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	} else if loginVia == nil || loginVia.ViaUsername() {
		domainUser, err = service.userRepository.GetByUsername(username)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	}

	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(domainUser.Password), []byte(password))
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	loginUser := service.userConverter.ToDto(domainUser)
	return loginUser, nil
}

func (service *AuthServiceImpl) GenerateAccessToken(user *user.UserDto) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(authconfig.Cfg.ExpiryDuration)).Unix()
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authconfig.Cfg.SecretKey))
}

func (service *AuthServiceImpl) GenerateRefreshToken(user *user.UserDto, expireLong bool) (string, error) {
	var expirationTime int64
	if expireLong {
		expirationTime = time.Now().Add(time.Second * time.Duration(authconfig.Cfg.ExpiryLongDuration)).Unix()
	} else {
		expirationTime = time.Now().Add(time.Second * time.Duration(authconfig.Cfg.ExpiryDuration)).Unix()
	}

	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authconfig.Cfg.SecretKey))
}

func (service *AuthServiceImpl) GetAllPermissions(user *user.UserDto) ([]*auth.PermissionDto, error) {
	allPermissions := make([]*auth.PermissionDto, 0)

	roles, err := service.roleRepository.GetRolesForUser(user.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, role := range roles {
		permissions, err := service.permissionRepository.GetAllByRoleId(role.ID)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
		for _, permission := range permissions {
			for _, name := range permission.GetAllActivePermissions() {
				allPermissions = append(allPermissions, &auth.PermissionDto{
					Module: permission.ModuleName,
					Name:   name,
				})
			}
		}
	}
	return allPermissions, nil
}
