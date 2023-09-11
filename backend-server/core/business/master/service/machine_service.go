package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"

type MachineService interface {
	GetAllMachines(page int16, pageSize int16, sort string, filter *machine.MachineFilterDto) ([]*machine.MachineDto, int64, error)
	CreateMachine(machineDto *machine.MachineCreateDto) error
	GetMachineByID(machineID int64) (*machine.MachineDto, error)
	GetMachineByCode(machineCode string) (*machine.MachineDto, error)
	UpdateMachine(machineID int64, machine *machine.MachineUpdateDto) error
	DeleteMachine(machineID int64) error
	DeleteMachineByIDs(machineIDs []int64) error
}
