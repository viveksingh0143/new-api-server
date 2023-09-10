package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
)

type MachineConverter struct{}

func NewMachineConverter() *MachineConverter {
	return &MachineConverter{}
}

func (c *MachineConverter) ToMinimalDto(domainMachine *domain.Machine) *machine.MachineMinimalDto {
	machineDto := &machine.MachineMinimalDto{
		ID:     domainMachine.ID,
		Code:   domainMachine.Code,
		Name:   domainMachine.Name,
		Status: domainMachine.Status,
	}
	return machineDto
}

func (c *MachineConverter) ToDto(domainMachine *domain.Machine) *machine.MachineDto {
	machineDto := &machine.MachineDto{
		ID:            domainMachine.ID,
		Code:          domainMachine.Code,
		Name:          domainMachine.Name,
		Status:        domainMachine.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainMachine.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainMachine.UpdatedAt),
		LastUpdatedBy: domainMachine.LastUpdatedBy,
	}
	return machineDto
}

func (c *MachineConverter) ToDtoSlice(domainMachines []*domain.Machine) []*machine.MachineDto {
	var machineDtos []*machine.MachineDto
	for _, domainMachine := range domainMachines {
		machineDtos = append(machineDtos, c.ToDto(domainMachine))
	}
	return machineDtos
}

func (c *MachineConverter) ToDomain(machineDto *machine.MachineCreateDto) *domain.Machine {
	domainMachine := &domain.Machine{
		Code:          machineDto.Code,
		Name:          machineDto.Name.String,
		Status:        machineDto.Status,
		LastUpdatedBy: machineDto.LastUpdatedBy,
	}
	return domainMachine
}

func (c *MachineConverter) ToUpdateDomain(domainMachine *domain.Machine, machineDto *machine.MachineUpdateDto) {
	if machineDto.Code != "" {
		domainMachine.Code = machineDto.Code
	}
	if machineDto.Name != "" {
		domainMachine.Name = machineDto.Name
	}
	if machineDto.Status.IsValid() {
		domainMachine.Status = machineDto.Status
	}
	domainMachine.LastUpdatedBy = machineDto.LastUpdatedBy
}
