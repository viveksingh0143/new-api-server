package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type ContainerServiceImpl struct {
	ContainerRepo      repository.ContainerRepository
	ProductRepo        repository.ProductRepository
	ContainerConverter *converter.ContainerConverter
}

func NewContainerService(containerRepo repository.ContainerRepository, productRepo repository.ProductRepository, containerConverter *converter.ContainerConverter) ContainerService {
	return &ContainerServiceImpl{ContainerRepo: containerRepo, ProductRepo: productRepo, ContainerConverter: containerConverter}
}

func (s *ContainerServiceImpl) GetAllContainers(page int16, pageSize int16, sort string, filter *container.ContainerFilterDto) ([]*container.ContainerDto, int64, error) {
	totalCount, err := s.ContainerRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainContainers, err := s.ContainerRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	// Convert domain containers to DTOs. You can do this based on your requirements.
	var containerDtos []*container.ContainerDto = s.ContainerConverter.ToDtoSlice(domainContainers)

	if len(containerDtos) > 0 {
		resourceIds := make([]int64, 0, len(containerDtos))
		for _, containerDto := range containerDtos {
			if containerDto.ResourceID.Valid && containerDto.ResourceName.Valid && containerDto.ResourceName.String == "*domain.Product" {
				resourceIds = append(resourceIds, containerDto.ResourceID.Int64)
			}
		}
		if len(resourceIds) > 0 {
			products, err := s.ProductRepo.GetByIds(resourceIds)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, 0, err
			}
			if len(products) > 0 {
				for _, containerDto := range containerDtos {
					if containerDto.ResourceID.Valid {
						for _, product := range products {
							if product.ID == containerDto.ResourceID.Int64 {
								containerDto.ContainerItemDto = &container.ContainerItemDto{
									ID:    product.ID,
									Code:  product.Code,
									Name:  product.Name,
									Type:  string(product.ProductType),
									Count: containerDto.ItemsCount,
								}
							}
						}
					}
				}
			}
		}
	}
	if len(containerDtos) > 0 {
		resourceIds := make([]int64, 0, len(containerDtos))
		for _, containerDto := range containerDtos {
			if containerDto.ResourceID.Valid && containerDto.ResourceName.Valid && containerDto.ResourceName.String == "*domain.Container" {
				resourceIds = append(resourceIds, containerDto.ResourceID.Int64)
			}
		}
		if len(resourceIds) > 0 {
			containerDomains, err := s.ContainerRepo.GetByIds(resourceIds)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, 0, err
			}
			if len(containerDomains) > 0 {
				for _, containerDto := range containerDtos {
					if containerDto.ResourceID.Valid {
						for _, containerDomain := range containerDomains {
							if containerDomain.ID == containerDto.ResourceID.Int64 {
								containerDto.ContainerItemDto = &container.ContainerItemDto{
									ID:    containerDomain.ID,
									Code:  containerDomain.Code,
									Name:  containerDomain.Name,
									Type:  string(containerDomain.ContainerType),
									Count: containerDto.ItemsCount,
								}
							}
						}
					}
				}
			}
		}
	}
	return containerDtos, int64(totalCount), nil
}

func (s *ContainerServiceImpl) CreateContainer(containerDto *container.ContainerCreateDto) error {
	var newContainer *domain.Container = s.ContainerConverter.ToDomain(containerDto)
	err := s.ContainerRepo.Create(newContainer)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ContainerServiceImpl) GetContainerByID(containerID int64) (*container.ContainerDto, error) {
	domainContainer, err := s.ContainerRepo.GetById(containerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	containerDto := s.ContainerConverter.ToDto(domainContainer)
	if containerDto.ResourceID.Valid && containerDto.ResourceName.Valid {
		if containerDto.ResourceName.String == "*domain.Product" {
			product, err := s.ProductRepo.GetById(containerDto.ResourceID.Int64)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, err
			}
			containerDto.ContainerItemDto = &container.ContainerItemDto{
				ID:    product.ID,
				Code:  product.Code,
				Name:  product.Name,
				Type:  string(product.ProductType),
				Count: containerDto.ItemsCount,
			}
		}
		if containerDto.ResourceName.String == "*domain.Container" {
			containerDomain, err := s.ContainerRepo.GetById(containerDto.ResourceID.Int64)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, err
			}
			containerDto.ContainerItemDto = &container.ContainerItemDto{
				ID:    containerDomain.ID,
				Code:  containerDomain.Code,
				Name:  containerDomain.Name,
				Type:  string(containerDomain.ContainerType),
				Count: containerDto.ItemsCount,
			}
		}
	}
	return containerDto, nil
}

func (s *ContainerServiceImpl) GetMinimalContainerByID(containerID int64) (*container.ContainerMinimalDto, error) {
	domainContainer, err := s.ContainerRepo.GetById(containerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ContainerConverter.ToMinimalDto(domainContainer), nil
}

func (s *ContainerServiceImpl) GetMinimalContainerByIds(containerIDs []int64) ([]*container.ContainerMinimalDto, error) {
	domainContainers, err := s.ContainerRepo.GetByIds(containerIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ContainerConverter.ToMinimalDtoSlice(domainContainers), nil
}

func (s *ContainerServiceImpl) GetContainerByCode(containerCode string) (*container.ContainerDto, error) {
	domainContainer, err := s.ContainerRepo.GetByCode(containerCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ContainerConverter.ToDto(domainContainer), nil
}

func (s *ContainerServiceImpl) UpdateContainer(containerID int64, containerDto *container.ContainerUpdateDto) error {
	existingContainer, err := s.ContainerRepo.GetById(containerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.ContainerConverter.ToUpdateDomain(existingContainer, containerDto)
	if err := s.ContainerRepo.Update(existingContainer); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ContainerServiceImpl) DeleteContainer(containerID int64) error {
	if err := s.ContainerRepo.Delete(containerID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ContainerServiceImpl) DeleteContainerByIDs(containerIDs []int64) error {
	if err := s.ContainerRepo.DeleteByIDs(containerIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ContainerServiceImpl) GetContainerCodeInfoDto() ([]*container.ContainerCodeInfoDto, error) {
	containersInfo, err := s.ContainerRepo.GetContainerCodeInfoDto()
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	existingMap := make(map[customtypes.ContainerType]bool)
	for _, container := range containersInfo {
		existingMap[container.ContainerType] = true
	}

	// Check for missing container types and add them with default code if necessary
	var finalContainers []*container.ContainerCodeInfoDto
	finalContainers = append(finalContainers, containersInfo...)

	for _, containerType := range customtypes.GetAllContainerTypes() {
		if _, exists := existingMap[containerType]; !exists {
			var defaultCode string = "UN00001"
			if containerType == customtypes.PALLET_TYPE {
				defaultCode = "PAL00001"
			} else if containerType == customtypes.BIN_TYPE {
				defaultCode = "BIN00001"
			} else if containerType == customtypes.RACK_TYPE {
				defaultCode = "RAC00001"
			}
			finalContainers = append(finalContainers, &container.ContainerCodeInfoDto{ContainerType: containerType, Code: defaultCode})
		}
	}

	return finalContainers, nil
}

func (s *ContainerServiceImpl) GetOneActiveContainerByCodeAndType(code string, containerType customtypes.ContainerType) (*container.ContainerMinimalDto, error) {
	domainContainer, err := s.ContainerRepo.GetOneContainerByCodeAndType(code, containerType)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ContainerConverter.ToMinimalDto(domainContainer), nil
}

func (s *ContainerServiceImpl) MarkContainerFullByCode(code string) error {
	domainContainer, err := s.ContainerRepo.GetByCode(code)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if !domainContainer.ResourceID.Valid || domainContainer.ResourceID.Int64 == 0 || domainContainer.Level.IsEmpty() {
		log.Printf("%+v\n", err)
		cType := strings.ToLower(domainContainer.ContainerType.String())
		return fmt.Errorf("%s (%s) is currently empty, first fill the %s", cType, code, cType)
	}

	if domainContainer.Level.IsFull() {
		log.Printf("%+v\n", err)
		cType := strings.ToLower(domainContainer.ContainerType.String())
		return fmt.Errorf("%s (%s) is already full", cType, code)
	}

	err = s.ContainerRepo.MarkContainerFullById(domainContainer.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ContainerServiceImpl) GetAllContainerApprovals(page int16, pageSize int16, sort string, filter *container.ContainerFilterDto) ([]*container.ContainerDto, int64, error) {
	filter.IsApproved = customtypes.NewValidNullBool(false)
	totalCount, err := s.ContainerRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainContainers, err := s.ContainerRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	// Convert domain containers to DTOs. You can do this based on your requirements.
	var containerDtos []*container.ContainerDto = s.ContainerConverter.ToDtoSlice(domainContainers)

	if len(containerDtos) > 0 {
		resourceIds := make([]int64, 0, len(containerDtos))
		for _, containerDto := range containerDtos {
			if containerDto.ResourceID.Valid && containerDto.ResourceName.Valid && containerDto.ResourceName.String == "*domain.Product" {
				resourceIds = append(resourceIds, containerDto.ResourceID.Int64)
			}
		}
		if len(resourceIds) > 0 {
			products, err := s.ProductRepo.GetByIds(resourceIds)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, 0, err
			}
			if len(products) > 0 {
				for _, containerDto := range containerDtos {
					if containerDto.ResourceID.Valid {
						for _, product := range products {
							if product.ID == containerDto.ResourceID.Int64 {
								containerDto.ContainerItemDto = &container.ContainerItemDto{
									ID:    product.ID,
									Code:  product.Code,
									Name:  product.Name,
									Type:  string(product.ProductType),
									Count: containerDto.ItemsCount,
								}
							}
						}
					}
				}
			}
		}
	}
	if len(containerDtos) > 0 {
		resourceIds := make([]int64, 0, len(containerDtos))
		for _, containerDto := range containerDtos {
			if containerDto.ResourceID.Valid && containerDto.ResourceName.Valid && containerDto.ResourceName.String == "*domain.Container" {
				resourceIds = append(resourceIds, containerDto.ResourceID.Int64)
			}
		}
		if len(resourceIds) > 0 {
			containerDomains, err := s.ContainerRepo.GetByIds(resourceIds)
			if err != nil {
				log.Printf("%+v\n", err)
				return nil, 0, err
			}
			if len(containerDomains) > 0 {
				for _, containerDto := range containerDtos {
					if containerDto.ResourceID.Valid {
						for _, containerDomain := range containerDomains {
							if containerDomain.ID == containerDto.ResourceID.Int64 {
								containerDto.ContainerItemDto = &container.ContainerItemDto{
									ID:    containerDomain.ID,
									Code:  containerDomain.Code,
									Name:  containerDomain.Name,
									Type:  string(containerDomain.ContainerType),
									Count: containerDto.ItemsCount,
								}
							}
						}
					}
				}
			}
		}
	}
	return containerDtos, int64(totalCount), nil
}

func (s *ContainerServiceImpl) ApproveContainerByIDs(containerIDs []int64) error {
	if err := s.ContainerRepo.ApproveContainerByIDs(containerIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
