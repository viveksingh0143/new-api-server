package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/requisition"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
)

type RequisitionService interface {
	GetAllRequisitions(page int16, pageSize int16, sort string, filter *requisition.RequisitionFilterDto) ([]*requisition.RequisitionDto, int64, error)
	CreateRequisition(requisitionDto *requisition.RequisitionCreateDto) error
	GetRequisitionByID(requisitionID int64) (*requisition.RequisitionDto, error)
	GetMinimalRequisitionByID(requisitionID int64) (*requisition.RequisitionMinimalDto, error)
	GetRequisitionByCode(requisitionCode string) (*requisition.RequisitionDto, []*reports.InventoryStatusDetail, error)
	UpdateRequisition(requisitionID int64, requisition *requisition.RequisitionUpdateDto) error
	DeleteRequisition(requisitionID int64) error
	DeleteRequisitionByIDs(requisitionIDs []int64) error
}
