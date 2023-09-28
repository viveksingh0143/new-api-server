package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
)

type ContainerConverter struct {
	storeConv StoreConverter
}

func NewContainerConverter(storeConverter StoreConverter) *ContainerConverter {
	return &ContainerConverter{
		storeConv: storeConverter,
	}
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
		IsApproved:    domainContainer.IsApproved,
		Status:        domainContainer.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainContainer.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainContainer.UpdatedAt),
		LastUpdatedBy: domainContainer.LastUpdatedBy,
		Level:         domainContainer.Level,
		StoreID:       domainContainer.StoreID,
		Store:         c.storeConv.ToMinimalDto(domainContainer.Store),
		OtherInfo:     domainContainer.Info(),
		ResourceID:    domainContainer.ResourceID,
		ResourceName:  domainContainer.ResourceName,
		ItemsCount:    domainContainer.ItemsCount,
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
	if containerDto.Store != nil && containerDto.Store.ID > 0 {
		domainContainer.StoreID = customtypes.NewValidNullInt64(containerDto.Store.ID)
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
	if containerDto.Store != nil && containerDto.Store.ID > 0 {
		domainContainer.StoreID = customtypes.NewValidNullInt64(containerDto.Store.ID)
	} else {
		domainContainer.StoreID = customtypes.NewInvalidNullInt64()
	}
}
