package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserConverter struct {
	roleConverter *RoleConverter
}

func NewUserConverter(roleConv *RoleConverter) *UserConverter {
	return &UserConverter{roleConverter: roleConv}
}

func (c *UserConverter) ToMinimalDto(domainUser *domain.User) user.UserMinimalDto {
	userDto := user.UserMinimalDto{
		ID:      domainUser.ID,
		Name:    domainUser.Name,
		StaffID: customtypes.NewValidNullString(domainUser.StaffID),
		EMail:   customtypes.NewValidNullString(domainUser.EMail),
		Status:  domainUser.Status,
	}
	return userDto
}

func (c *UserConverter) ToDto(domainUser *domain.User) user.UserDto {
	userDto := user.UserDto{
		ID:            domainUser.ID,
		Name:          domainUser.Name,
		Username:      customtypes.NewValidNullString(domainUser.Username),
		StaffID:       customtypes.NewValidNullString(domainUser.StaffID),
		EMail:         customtypes.NewValidNullString(domainUser.EMail),
		Status:        domainUser.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainUser.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainUser.UpdatedAt),
		LastUpdatedBy: domainUser.LastUpdatedBy,
	}
	if domainUser.EMailConfirmationAt != nil {
		userDto.EMailConfirmed = true
	} else {
		userDto.EMailConfirmed = false
	}

	// Convert Roles from domain model to DTO
	var roleDtos []role.RoleMinimalDto
	for _, domainRole := range domainUser.Roles {
		roleDtos = append(roleDtos, c.roleConverter.ToMinimalDto(domainRole))
	}
	userDto.Roles = roleDtos

	return userDto
}

func (c *UserConverter) ToDtoSlice(domainUsers []*domain.User) []user.UserDto {
	var userDtos []user.UserDto
	for _, domainUser := range domainUsers {
		userDtos = append(userDtos, c.ToDto(domainUser))
	}
	return userDtos
}

func (c *UserConverter) ToDomain(userDto user.UserCreateDto) *domain.User {
	domainUser := &domain.User{
		Name:          userDto.Name,
		Username:      userDto.Username.String,
		StaffID:       userDto.StaffID.String,
		EMail:         userDto.EMail.String,
		Status:        userDto.Status,
		LastUpdatedBy: userDto.LastUpdatedBy,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err == nil {
		domainUser.Password = string(hashedPassword)
	}

	var domainRoles []*domain.Role
	for _, roleDto := range userDto.Roles {
		domainRole := &domain.Role{ID: roleDto.ID}
		domainRoles = append(domainRoles, domainRole)
	}
	domainUser.Roles = domainRoles
	return domainUser
}

func (c *UserConverter) ToUpdateDomain(domainUser *domain.User, userDto user.UserUpdateDto) {

	if userDto.Name != "" {
		domainUser.Name = userDto.Name
	}
	if userDto.Username.Valid {
		domainUser.Username = userDto.Username.String
	}
	if userDto.StaffID.Valid {
		domainUser.StaffID = userDto.StaffID.String
	}
	if userDto.EMail.Valid {
		domainUser.EMail = userDto.EMail.String
	}
	if userDto.Status.IsValid() {
		domainUser.Status = userDto.Status
	}
	domainUser.LastUpdatedBy = userDto.LastUpdatedBy
	if userDto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
		if err == nil {
			domainUser.Password = string(hashedPassword)
		}
	}

	var domainRoles []*domain.Role
	for _, roleDto := range userDto.Roles {
		domainRole := &domain.Role{ID: roleDto.ID}
		domainRoles = append(domainRoles, domainRole)
	}
	domainUser.Roles = domainRoles
}
