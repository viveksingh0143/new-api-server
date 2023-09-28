package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/core/business/master/reports"
)

type OutwardRequestRepository interface {
	GetTotalCount(filter *outwardrequest.OutwardRequestFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *outwardrequest.OutwardRequestFilterDto) ([]*domain.OutwardRequest, error)
	Create(outwardrequest *domain.OutwardRequest) error
	GetById(outwardrequestID int64) (*domain.OutwardRequest, error)
	GetByOrderNo(outwardrequestOrderNo string) (*domain.OutwardRequest, error)
	Update(outwardrequest *domain.OutwardRequest) error
	Delete(outwardrequestID int64) error
	DeleteByIDs(outwardrequestIDs []int64) error
	GetItemsForOutwardRequests(orderIDs []int64) (map[int64][]*domain.OutwardRequestItem, error)
	GetItemsForOutwardRequest(orderIDs int64) ([]*domain.OutwardRequestItem, error)
	GetShipperLabels(requestID int64, requestName string) ([]*reports.OutwardRequestShipperReport, error)
	GenerateShipperLabels(shipperLabel *domain.ShipperLabel, batchNo string, productID int64) error
}
