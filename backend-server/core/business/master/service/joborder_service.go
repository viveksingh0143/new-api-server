package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"

type JobOrderService interface {
	GetAllJobOrders(page int16, pageSize int16, sort string, filter *joborder.JobOrderFilterDto) ([]*joborder.JobOrderDto, int64, error)
	CreateJobOrder(joborderDto *joborder.JobOrderCreateDto) error
	GetJobOrderByID(joborderID int64) (*joborder.JobOrderDto, error)
	GetMinimalJobOrderByID(joborderID int64) (*joborder.JobOrderMinimalDto, error)
	GetJobOrderByCode(joborderCode string) (*joborder.JobOrderDto, error)
	UpdateJobOrder(joborderID int64, joborder *joborder.JobOrderUpdateDto) error
	DeleteJobOrder(joborderID int64) error
	DeleteJobOrderByIDs(joborderIDs []int64) error
}
