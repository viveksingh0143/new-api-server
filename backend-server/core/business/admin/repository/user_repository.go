package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(userID int64) error
	DeleteByIDs(userIDs []int64) error
	GetById(userID int64) (*domain.User, error)
	GetByIds(userIDs []int64) ([]*domain.User, error)
	GetTotalCount(filter *user.UserFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *user.UserFilterDto) ([]*domain.User, error)

	GetByUsername(username string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByStaffID(staffID string) (*domain.User, error)
}
