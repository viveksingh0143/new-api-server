package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ContainerType string

const (
	PALLET_TYPE ContainerType = "PALLET"
	BIN_TYPE    ContainerType = "BIN"
	RACK_TYPE   ContainerType = "RACK"
)

func GetAllContainerTypes() []ContainerType {
	return []ContainerType{PALLET_TYPE, BIN_TYPE, RACK_TYPE}
}

func (s ContainerType) IsValid() bool {
	return s == PALLET_TYPE || s == BIN_TYPE || s == RACK_TYPE
}

func (s ContainerType) String() string {
	switch s {
	case PALLET_TYPE:
		return "PALLET"
	case BIN_TYPE:
		return "BIN"
	case RACK_TYPE:
		return "RACK"
	default:
		return "UNKNOWN"
	}
}

func ValidateContainerType(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(ContainerType)
	if !ok {
		return false
	}
	value := ContainerType(strings.ToUpper(string(rawValue)))

	switch value {
	case PALLET_TYPE, BIN_TYPE, RACK_TYPE:
		return true
	default:
		return false
	}
}

// MarshalJSON for ContainerType
func (s ContainerType) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid container type: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for ContainerType
func (s *ContainerType) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "PALLET":
		*s = PALLET_TYPE
	case "BIN":
		*s = BIN_TYPE
	case "RACK":
		*s = RACK_TYPE
	default:
		return fmt.Errorf("invalid container type: %s", str)
	}
	return nil
}
