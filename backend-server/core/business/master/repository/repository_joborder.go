package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"
)

type JobOrderRepository interface {
	GetTotalCount(filter *joborder.JobOrderFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *joborder.JobOrderFilterDto) ([]*domain.JobOrder, error)
	Create(joborder *domain.JobOrder) error
	GetById(joborderID int64) (*domain.JobOrder, error)
	GetByOrderNo(joborderOrderNo string) (*domain.JobOrder, error)
	Update(joborder *domain.JobOrder) error
	Delete(joborderID int64) error
	DeleteByIDs(joborderIDs []int64) error
	GetItemsForJobOrders(orderIDs []int64) (map[int64][]*domain.JobOrderItem, error)
	GetItemsForJobOrder(orderIDs int64) ([]*domain.JobOrderItem, error)
}
