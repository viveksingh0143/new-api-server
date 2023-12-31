package repository

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
)

type ContainerRepository interface {
	Create(container *domain.Container) error
	Update(container *domain.Container) error
	Delete(containerID int64) error
	DeleteByIDs(containerIDs []int64) error
	GetById(containerID int64) (*domain.Container, error)
	GetByIds(containerIDs []int64) ([]*domain.Container, error)
	GetByCode(containerCode string) (*domain.Container, error)
	GetTotalCount(filter *container.ContainerFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *container.ContainerFilterDto) ([]*domain.Container, error)
	GetContainerCodeInfoDto() ([]*container.ContainerCodeInfoDto, error)
	GetOneContainerByCodeAndType(code string, containerType customtypes.ContainerType) (*domain.Container, error)
	MarkContainerFullById(id int64) error
	AttachedCount(resourceId int64, resourceType string) (int, error)

	ApproveContainerByIDs(containerIDs []int64) error
}
