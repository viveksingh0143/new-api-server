package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type JobOrderServiceImpl struct {
	JobOrderRepo      repository.JobOrderRepository
	CustomerRepo      repository.CustomerRepository
	ProductRepo       repository.ProductRepository
	JobOrderConverter *converter.JobOrderConverter
}

func NewJobOrderService(jobOrderRepo repository.JobOrderRepository, customerRepo repository.CustomerRepository, productRepo repository.ProductRepository, jobOrderConverter *converter.JobOrderConverter) JobOrderService {
	return &JobOrderServiceImpl{JobOrderRepo: jobOrderRepo, CustomerRepo: customerRepo, ProductRepo: productRepo, JobOrderConverter: jobOrderConverter}
}

func (s *JobOrderServiceImpl) GetAllJobOrders(page int16, pageSize int16, sort string, filter *joborder.JobOrderFilterDto) ([]*joborder.JobOrderDto, int64, error) {
	totalCount, err := s.JobOrderRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainJobOrders, err := s.JobOrderRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}

	if len(domainJobOrders) > 0 {
		jobOrderIds := make([]int64, 0, len(domainJobOrders))
		for _, jobOrder := range domainJobOrders {
			jobOrderIds = append(jobOrderIds, jobOrder.ID)
		}

		jobOrderItemsMap, err := s.JobOrderRepo.GetItemsForJobOrders(jobOrderIds)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}
		for i, jobOrder := range domainJobOrders {
			if roles, ok := jobOrderItemsMap[jobOrder.ID]; ok {
				domainJobOrders[i].Items = roles
			}
		}
	}

	// Convert domain joborders to DTOs. You can do this based on your requirements.
	var joborderDtos []*joborder.JobOrderDto = s.JobOrderConverter.ToDtoSlice(domainJobOrders)
	return joborderDtos, int64(totalCount), nil
}

func (s *JobOrderServiceImpl) CreateJobOrder(joborderDto *joborder.JobOrderCreateDto) error {
	var newJobOrder *domain.JobOrder = s.JobOrderConverter.ToDomain(joborderDto)
	err := s.JobOrderRepo.Create(newJobOrder)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *JobOrderServiceImpl) GetJobOrderByID(joborderID int64) (*joborder.JobOrderDto, error) {
	domainJobOrder, err := s.JobOrderRepo.GetById(joborderID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	domainCustomer, err := s.CustomerRepo.GetById(domainJobOrder.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainJobOrder.Customer = domainCustomer

	items, err := s.JobOrderRepo.GetItemsForJobOrder(domainJobOrder.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, item := range items {
		productDomain, err := s.ProductRepo.GetById(item.ProductID)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
		item.Product = productDomain
	}
	domainJobOrder.Items = items

	return s.JobOrderConverter.ToDto(domainJobOrder), nil
}

func (s *JobOrderServiceImpl) GetMinimalJobOrderByID(joborderID int64) (*joborder.JobOrderMinimalDto, error) {
	domainJobOrder, err := s.JobOrderRepo.GetById(joborderID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.JobOrderConverter.ToMinimalDto(domainJobOrder), nil
}

func (s *JobOrderServiceImpl) GetJobOrderByCode(joborderCode string) (*joborder.JobOrderDto, error) {
	domainJobOrder, err := s.JobOrderRepo.GetByOrderNo(joborderCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.JobOrderConverter.ToDto(domainJobOrder), nil
}

func (s *JobOrderServiceImpl) UpdateJobOrder(joborderID int64, joborderDto *joborder.JobOrderUpdateDto) error {
	existingJobOrder, err := s.JobOrderRepo.GetById(joborderID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.JobOrderConverter.ToUpdateDomain(existingJobOrder, joborderDto)
	if err := s.JobOrderRepo.Update(existingJobOrder); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *JobOrderServiceImpl) DeleteJobOrder(joborderID int64) error {
	if err := s.JobOrderRepo.Delete(joborderID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *JobOrderServiceImpl) DeleteJobOrderByIDs(joborderIDs []int64) error {
	if err := s.JobOrderRepo.DeleteByIDs(joborderIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
