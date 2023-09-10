package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type MachineServiceImpl struct {
	MachineRepo      repository.MachineRepository
	MachineConverter *converter.MachineConverter
}

func NewMachineService(machineRepo repository.MachineRepository, machineConverter *converter.MachineConverter) MachineService {
	return &MachineServiceImpl{MachineRepo: machineRepo, MachineConverter: machineConverter}
}

func (s *MachineServiceImpl) GetAllMachines(page int64, pageSize int64, sort string, filter *machine.MachineFilterDto) ([]*machine.MachineDto, int64, error) {
	totalCount, err := s.MachineRepo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}
	domainMachines, err := s.MachineRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		return nil, 0, err
	}
	// Convert domain machines to DTOs. You can do this based on your requirements.
	var machineDtos []*machine.MachineDto = s.MachineConverter.ToDtoSlice(domainMachines)
	return machineDtos, int64(totalCount), nil
}

func (s *MachineServiceImpl) CreateMachine(machineDto *machine.MachineCreateDto) error {
	var newMachine *domain.Machine = s.MachineConverter.ToDomain(machineDto)
	err := s.MachineRepo.Create(newMachine)
	if err != nil {
		return err
	}
	return nil
}

func (s *MachineServiceImpl) GetMachineByID(machineID int64) (*machine.MachineDto, error) {
	domainMachine, err := s.MachineRepo.GetById(machineID)
	if err != nil {
		return nil, err
	}
	return s.MachineConverter.ToDto(domainMachine), nil
}

func (s *MachineServiceImpl) GetMachineByCode(machineCode string) (*machine.MachineDto, error) {
	domainMachine, err := s.MachineRepo.GetByCode(machineCode)
	if err != nil {
		return nil, err
	}
	return s.MachineConverter.ToDto(domainMachine), nil
}

func (s *MachineServiceImpl) UpdateMachine(machineID int64, machineDto *machine.MachineUpdateDto) error {
	existingMachine, err := s.MachineRepo.GetById(machineID)
	if err != nil {
		return err
	}

	s.MachineConverter.ToUpdateDomain(existingMachine, machineDto)
	if err := s.MachineRepo.Update(existingMachine); err != nil {
		return err
	}
	return nil
}

func (s *MachineServiceImpl) DeleteMachine(machineID int64) error {
	if err := s.MachineRepo.Delete(machineID); err != nil {
		return err
	}
	return nil
}

func (s *MachineServiceImpl) DeleteMachineByIDs(machineIDs []int64) error {
	if err := s.MachineRepo.DeleteByIDs(machineIDs); err != nil {
		return err
	}
	return nil
}
