package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(userID int64) error
	GetById(userID int64) (*domain.User, error)
	GetByName(userName string) (*domain.User, error)
	GetTotalCount(filter user.UserFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter user.UserFilterDto) ([]*domain.User, error)
}
