package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/permission"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
)

type RoleConverter struct{}

func NewRoleConverter() *RoleConverter {
	return &RoleConverter{}
}

func (c *RoleConverter) ToMinimalDto(domainRole *domain.Role) *role.RoleMinimalDto {
	roleDto := &role.RoleMinimalDto{
		ID:     domainRole.ID,
		Name:   domainRole.Name,
		Status: domainRole.Status,
	}
	return roleDto
}

func (c *RoleConverter) ToDto(domainRole *domain.Role) *role.RoleDto {
	roleDto := &role.RoleDto{
		ID:            domainRole.ID,
		Name:          domainRole.Name,
		Status:        domainRole.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainRole.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainRole.UpdatedAt),
		LastUpdatedBy: domainRole.LastUpdatedBy,
	}

	// Convert Permissions from domain model to DTO
	var permissionDtos []*permission.PermissionDto = make([]*permission.PermissionDto, 0)
	for _, domainPermission := range domainRole.Permissions {
		permissionDto := &permission.PermissionDto{
			ID:         domainPermission.ID,
			RoleID:     domainPermission.RoleID,
			ModuleName: domainPermission.ModuleName,
			Create:     domainPermission.Create,
			Read:       domainPermission.Read,
			Update:     domainPermission.Update,
			Delete:     domainPermission.Delete,
			Export:     domainPermission.Export,
			Import:     domainPermission.Import,
		}
		permissionDtos = append(permissionDtos, permissionDto)
	}
	roleDto.Permissions = permissionDtos

	return roleDto
}

func (c *RoleConverter) ToDtoSlice(domainRoles []*domain.Role) []*role.RoleDto {
	var roleDtos = make([]*role.RoleDto, 0)
	for _, domainRole := range domainRoles {
		roleDtos = append(roleDtos, c.ToDto(domainRole))
	}
	return roleDtos
}

func (c *RoleConverter) ToDomain(roleDto *role.RoleCreateDto) *domain.Role {
	domainRole := &domain.Role{
		Name:          roleDto.Name.String,
		Status:        roleDto.Status.Status,
		LastUpdatedBy: roleDto.LastUpdatedBy,
	}

	var domainPermissions []*domain.Permission
	for _, permissionDto := range roleDto.Permissions {
		domainPermission := c.ToPermissionDomain(permissionDto)
		domainPermissions = append(domainPermissions, domainPermission)
	}
	domainRole.Permissions = domainPermissions

	return domainRole
}

func (c *RoleConverter) ToUpdateDomain(domainRole *domain.Role, roleDto *role.RoleUpdateDto) {
	if roleDto.Name != "" {
		domainRole.Name = roleDto.Name
	}
	if roleDto.Status.IsValid() {
		domainRole.Status = roleDto.Status
	}
	domainRole.LastUpdatedBy = roleDto.LastUpdatedBy

	var domainPermissions []*domain.Permission
	for _, permissionDto := range roleDto.Permissions {
		domainPermission := c.ToPermissionDomain(permissionDto)
		domainPermissions = append(domainPermissions, domainPermission)
	}
	domainRole.Permissions = domainPermissions
}

func (c *RoleConverter) ToPermissionDomain(permissionDto *permission.PermissionDto) *domain.Permission {
	domainPermission := &domain.Permission{
		ID:         permissionDto.ID,
		RoleID:     permissionDto.RoleID,
		ModuleName: permissionDto.ModuleName,
		Create:     permissionDto.Create,
		Read:       permissionDto.Read,
		Update:     permissionDto.Update,
		Delete:     permissionDto.Delete,
		Export:     permissionDto.Export,
		Import:     permissionDto.Import,
	}

	return domainPermission
}
