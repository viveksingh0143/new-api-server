package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
)

type OutwardRequestService interface {
	GetAllOutwardRequests(page int16, pageSize int16, sort string, filter *outwardrequest.OutwardRequestFilterDto) ([]*outwardrequest.OutwardRequestDto, int64, error)
	CreateOutwardRequest(outwardrequestDto *outwardrequest.OutwardRequestCreateDto) error
	GetOutwardRequestByID(outwardrequestID int64) (*outwardrequest.OutwardRequestDto, error)
	GetMinimalOutwardRequestByID(outwardrequestID int64) (*outwardrequest.OutwardRequestMinimalDto, error)
	GetOutwardRequestByCode(outwardrequestCode string) (*outwardrequest.OutwardRequestDto, []*reports.InventoryRackStatusDetail, []*reports.InventoryBinStatusDetail, error)
	UpdateOutwardRequest(outwardrequestID int64, outwardrequest *outwardrequest.OutwardRequestUpdateDto) error
	DeleteOutwardRequest(outwardrequestID int64) error
	DeleteOutwardRequestByIDs(outwardrequestIDs []int64) error
}
