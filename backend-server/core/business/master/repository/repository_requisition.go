package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/requisition"
)

type RequisitionRepository interface {
	GetTotalCount(filter *requisition.RequisitionFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *requisition.RequisitionFilterDto) ([]*domain.Requisition, error)
	Create(requisition *domain.Requisition) error
	GetById(requisitionID int64) (*domain.Requisition, error)
	GetByOrderNo(requisitionOrderNo string) (*domain.Requisition, error)
	Update(requisition *domain.Requisition) error
	Delete(requisitionID int64) error
	DeleteByIDs(requisitionIDs []int64) error
	GetItemsForRequisitions(orderIDs []int64) (map[int64][]*domain.RequisitionItem, error)
	GetItemsForRequisition(orderIDs int64) ([]*domain.RequisitionItem, error)
}
