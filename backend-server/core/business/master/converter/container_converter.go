package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
)

type ContainerConverter struct{}

func NewContainerConverter() *ContainerConverter {
	return &ContainerConverter{}
}

func (c *ContainerConverter) ToMinimalDto(domainContainer *domain.Container) *container.ContainerMinimalDto {
	containerDto := &container.ContainerMinimalDto{
		ID:            domainContainer.ID,
		ContainerType: domainContainer.ContainerType,
		Code:          domainContainer.Code,
		Name:          domainContainer.Name,
		Status:        domainContainer.Status,
	}
	return containerDto
}

func (c *ContainerConverter) ToDto(domainContainer *domain.Container) *container.ContainerDto {
	containerDto := &container.ContainerDto{
		ID:            domainContainer.ID,
		ContainerType: domainContainer.ContainerType,
		Code:          domainContainer.Code,
		Name:          domainContainer.Name,
		Address:       domainContainer.Address,
		Status:        domainContainer.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainContainer.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainContainer.UpdatedAt),
		LastUpdatedBy: domainContainer.LastUpdatedBy,
		Level:         domainContainer.Level,
	}
	if containerDto.ContainerType == customtypes.PALLET_TYPE {
		containerDto.MinCapacity = domain.PalletContainerInfo.MinCapacity
		containerDto.MaxCapacity = domain.PalletContainerInfo.MaxCapacity
		containerDto.CanContains = domain.PalletContainerInfo.CanContains
	} else if containerDto.ContainerType == customtypes.BIN_TYPE {
		containerDto.MinCapacity = domain.BinContainerInfo.MinCapacity
		containerDto.MaxCapacity = domain.BinContainerInfo.MaxCapacity
		containerDto.CanContains = domain.BinContainerInfo.CanContains
	} else if containerDto.ContainerType == customtypes.RACK_TYPE {
		containerDto.MinCapacity = domain.RackContainerInfo.MinCapacity
		containerDto.MaxCapacity = domain.RackContainerInfo.MaxCapacity
		containerDto.CanContains = domain.RackContainerInfo.CanContains
	}
	return containerDto
}

func (c *ContainerConverter) ToDtoSlice(domainContainers []*domain.Container) []*container.ContainerDto {
	var containerDtos = make([]*container.ContainerDto, 0)
	for _, domainContainer := range domainContainers {
		containerDtos = append(containerDtos, c.ToDto(domainContainer))
	}
	return containerDtos
}

func (c *ContainerConverter) ToMinimalDtoSlice(domainContainers []*domain.Container) []*container.ContainerMinimalDto {
	var containerDtos = make([]*container.ContainerMinimalDto, 0)
	for _, domainContainer := range domainContainers {
		containerDtos = append(containerDtos, c.ToMinimalDto(domainContainer))
	}
	return containerDtos
}

func (c *ContainerConverter) ToDomain(containerDto *container.ContainerCreateDto) *domain.Container {
	domainContainer := &domain.Container{
		ContainerType: containerDto.ContainerType,
		Code:          containerDto.Code,
		Name:          containerDto.Name.String,
		Address:       containerDto.Address,
		Status:        containerDto.Status,
		LastUpdatedBy: containerDto.LastUpdatedBy,
	}
	return domainContainer
}

func (c *ContainerConverter) ToUpdateDomain(domainContainer *domain.Container, containerDto *container.ContainerUpdateDto) {
	domainContainer.ContainerType = containerDto.ContainerType
	domainContainer.Code = containerDto.Code
	domainContainer.Name = containerDto.Name.String
	domainContainer.Address = containerDto.Address
	if containerDto.Status.IsValid() {
		domainContainer.Status = containerDto.Status
	}
	domainContainer.LastUpdatedBy = containerDto.LastUpdatedBy
}
