package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
	"github.com/vamika-digital/wms-api-server/core/business/admin/repository"
)

type UserServiceImpl struct {
	UserRepo      repository.UserRepository
	RoleRepo      repository.RoleRepository
	UserConverter *converter.UserConverter
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, userConverter *converter.UserConverter) UserService {
	return &UserServiceImpl{UserRepo: userRepo, RoleRepo: roleRepo, UserConverter: userConverter}
}

func (s *UserServiceImpl) GetAllUsers(page int64, pageSize int64, sort string, filter *user.UserFilterDto) ([]*user.UserDto, int64, error) {
	totalCount, err := s.UserRepo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}
	domainUsers, err := s.UserRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(domainUsers) > 0 {
		userIds := make([]int64, 0, len(domainUsers))
		for _, user := range domainUsers {
			userIds = append(userIds, user.ID)
		}

		userRolesMap, err := s.RoleRepo.GetRolesForUsers(userIds)
		if err != nil {
			return nil, 0, err
		}
		for i, user := range domainUsers {
			if roles, ok := userRolesMap[user.ID]; ok {
				domainUsers[i].Roles = roles
			}
		}
	}

	var userDtos []*user.UserDto = s.UserConverter.ToDtoSlice(domainUsers)
	return userDtos, int64(totalCount), nil
}

func (s *UserServiceImpl) CreateUser(userDto *user.UserCreateDto) error {
	var newUser *domain.User = s.UserConverter.ToDomain(userDto)
	if err := s.UserRepo.Create(newUser); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) GetUserByID(userID int64) (*user.UserDto, error) {
	domainUser, err := s.UserRepo.GetById(userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepo.GetRolesForUser(domainUser.ID)
	if err != nil {
		return nil, err
	}
	domainUser.Roles = roles
	return s.UserConverter.ToDto(domainUser), nil
}

func (s *UserServiceImpl) UpdateUser(userID int64, userDto *user.UserUpdateDto) error {
	existingUser, err := s.UserRepo.GetById(userID)
	if err != nil {
		return err
	}

	s.UserConverter.ToUpdateDomain(existingUser, userDto)
	if err := s.UserRepo.Update(existingUser); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(userID int64) error {
	if err := s.UserRepo.Delete(userID); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) GetByUsername(username string) (*user.UserDto, error) {
	domainUser, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepo.GetRolesForUser(domainUser.ID)
	if err != nil {
		return nil, err
	}
	domainUser.Roles = roles
	return s.UserConverter.ToDto(domainUser), nil
}

func (s *UserServiceImpl) GetByEmail(email string) (*user.UserDto, error) {
	domainUser, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepo.GetRolesForUser(domainUser.ID)
	if err != nil {
		return nil, err
	}
	domainUser.Roles = roles
	return s.UserConverter.ToDto(domainUser), nil
}

func (s *UserServiceImpl) GetByStaffID(staffID string) (*user.UserDto, error) {
	domainUser, err := s.UserRepo.GetByStaffID(staffID)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepo.GetRolesForUser(domainUser.ID)
	if err != nil {
		return nil, err
	}
	domainUser.Roles = roles
	return s.UserConverter.ToDto(domainUser), nil
}

func (s *UserServiceImpl) DeleteUserByIDs(userIDs []int64) error {
	if err := s.UserRepo.DeleteByIDs(userIDs); err != nil {
		return err
	}
	return nil
}
