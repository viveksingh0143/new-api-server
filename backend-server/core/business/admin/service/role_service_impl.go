package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
	"github.com/vamika-digital/wms-api-server/core/business/admin/repository"
)

type RoleServiceImpl struct {
	RoleRepo       repository.RoleRepository
	PermissionRepo repository.PermissionRepository
	RoleConverter  *converter.RoleConverter
}

func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository, roleConverter *converter.RoleConverter) RoleService {
	return &RoleServiceImpl{RoleRepo: roleRepo, PermissionRepo: permissionRepo, RoleConverter: roleConverter}
}

func (s *RoleServiceImpl) GetAllRoles(page int64, pageSize int64, sort string, filter *role.RoleFilterDto) ([]*role.RoleDto, int64, error) {
	totalCount, err := s.RoleRepo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}
	domainRoles, err := s.RoleRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		return nil, 0, err
	}
	// Convert domain roles to DTOs. You can do this based on your requirements.
	var roleDtos []*role.RoleDto = s.RoleConverter.ToDtoSlice(domainRoles)
	return roleDtos, int64(totalCount), nil
}

func (s *RoleServiceImpl) CreateRole(roleDto *role.RoleCreateDto) error {
	var newRole *domain.Role = s.RoleConverter.ToDomain(roleDto)
	err := s.RoleRepo.Create(newRole)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoleServiceImpl) GetRoleByID(roleID int64) (*role.RoleDto, error) {
	domainRole, err := s.RoleRepo.GetById(roleID)
	if err != nil {
		return nil, err
	}
	permissions, err := s.PermissionRepo.GetAllByRoleId(domainRole.ID)
	if err != nil {
		return nil, err
	}
	domainRole.Permissions = permissions
	return s.RoleConverter.ToDto(domainRole), nil
}

func (s *RoleServiceImpl) UpdateRole(roleID int64, roleDto *role.RoleUpdateDto) error {
	existingRole, err := s.RoleRepo.GetById(roleID)
	if err != nil {
		return err
	}

	s.RoleConverter.ToUpdateDomain(existingRole, roleDto)
	if err := s.RoleRepo.Update(existingRole); err != nil {
		return err
	}
	return nil
}

func (s *RoleServiceImpl) DeleteRole(roleID int64) error {
	if err := s.RoleRepo.Delete(roleID); err != nil {
		return err
	}
	return nil
}

func (s *RoleServiceImpl) DeleteRoleByIDs(roleIDs []int64) error {
	if err := s.RoleRepo.DeleteByIDs(roleIDs); err != nil {
		return err
	}
	return nil
}
