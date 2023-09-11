package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ContainableType string

const (
	PALLET_CONTAINABLE        ContainableType = "PALLET_CONTAINABLE"
	BIN_CONTAINABLE           ContainableType = "BIN_CONTAINABLE"
	RACK_CONTAINABLE          ContainableType = "RACK_CONTAINABLE"
	RAW_MATERIAL_CONTAINABLE  ContainableType = "RAW_MATERIAL_CONTAINABLE"
	FINISHED_GOOD_CONTAINABLE ContainableType = "FINISHED_GOOD_CONTAINABLE"
)

func (s ContainableType) String() string {
	switch s {
	case PALLET_CONTAINABLE:
		return "PALLET_CONTAINABLE"
	case BIN_CONTAINABLE:
		return "BIN_CONTAINABLE"
	case RACK_CONTAINABLE:
		return "RACK_CONTAINABLE"
	case RAW_MATERIAL_CONTAINABLE:
		return "RAW_MATERIAL_CONTAINABLE"
	case FINISHED_GOOD_CONTAINABLE:
		return "FINISHED_GOOD_CONTAINABLE"
	default:
		return "UNKNOWN"
	}
}

func ValidateContainable(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(ContainableType)
	if !ok {
		return false
	}
	value := ContainableType(strings.ToUpper(string(rawValue)))

	switch value {
	case PALLET_CONTAINABLE, BIN_CONTAINABLE, RACK_CONTAINABLE, RAW_MATERIAL_CONTAINABLE, FINISHED_GOOD_CONTAINABLE:
		return true
	default:
		return false
	}
}

// MarshalJSON for ContainableType
func (s ContainableType) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid containable type: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for ContainableType
func (s *ContainableType) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "PALLET_CONTAINABLE":
		*s = PALLET_CONTAINABLE
	case "BIN_CONTAINABLE":
		*s = BIN_CONTAINABLE
	case "RACK_CONTAINABLE":
		*s = RACK_CONTAINABLE
	case "RAW_MATERIAL_CONTAINABLE":
		*s = RAW_MATERIAL_CONTAINABLE
	case "FINISHED_GOOD_CONTAINABLE":
		*s = FINISHED_GOOD_CONTAINABLE
	default:
		return fmt.Errorf("invalid containable type: %s", str)
	}
	return nil
}
