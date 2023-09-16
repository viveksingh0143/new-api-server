package service

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
)

type ContainerService interface {
	GetAllContainers(page int16, pageSize int16, sort string, filter *container.ContainerFilterDto) ([]*container.ContainerDto, int64, error)
	CreateContainer(containerDto *container.ContainerCreateDto) error
	GetContainerByID(containerID int64) (*container.ContainerDto, error)
	GetMinimalContainerByID(containerID int64) (*container.ContainerMinimalDto, error)
	GetMinimalContainerByIds(containerIDs []int64) ([]*container.ContainerMinimalDto, error)
	GetContainerByCode(containerCode string) (*container.ContainerDto, error)
	UpdateContainer(containerID int64, container *container.ContainerUpdateDto) error
	DeleteContainer(containerID int64) error
	DeleteContainerByIDs(containerIDs []int64) error
	GetContainerCodeInfoDto() ([]*container.ContainerCodeInfoDto, error)

	GetOneActiveContainerByCodeAndType(code string, containerType customtypes.ContainerType) (*container.ContainerMinimalDto, error)
	MarkContainerFullByCode(code string) error
}
