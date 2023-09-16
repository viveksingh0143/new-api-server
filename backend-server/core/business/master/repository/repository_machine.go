package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
)

type MachineRepository interface {
	Create(machine *domain.Machine) error
	Update(machine *domain.Machine) error
	Delete(machineID int64) error
	DeleteByIDs(machineIDs []int64) error
	GetById(machineID int64) (*domain.Machine, error)
	GetByIds(machineIDs []int64) ([]*domain.Machine, error)
	GetByCode(machineCode string) (*domain.Machine, error)
	GetByName(machineName string) (*domain.Machine, error)
	GetTotalCount(filter *machine.MachineFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *machine.MachineFilterDto) ([]*domain.Machine, error)
}
