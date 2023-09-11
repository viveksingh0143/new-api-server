package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UnitWeightType string

const (
	UnitWeightKilogram   UnitWeightType = "Kilogram"
	UnitWeightGram       UnitWeightType = "Gram"
	UnitWeightLiter      UnitWeightType = "Liter"
	UnitWeightMilliliter UnitWeightType = "Milliliter"
)

func GetAllUnitWeightTypes() []UnitWeightType {
	return []UnitWeightType{UnitWeightKilogram, UnitWeightGram, UnitWeightLiter, UnitWeightMilliliter}
}

func (s UnitWeightType) IsValid() bool {
	return s == UnitWeightKilogram || s == UnitWeightGram || s == UnitWeightLiter || s == UnitWeightMilliliter
}

func (s UnitWeightType) String() string {
	switch s {
	case UnitWeightKilogram:
		return "Kilogram"
	case UnitWeightGram:
		return "Gram"
	case UnitWeightLiter:
		return "Liter"
	case UnitWeightMilliliter:
		return "Milliliter"
	default:
		return "UNKNOWN"
	}
}

func ValidateUnitWeightType(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(UnitWeightType)
	if !ok {
		return false
	}
	value := UnitWeightType(string(rawValue))

	switch value {
	case UnitWeightKilogram, UnitWeightGram, UnitWeightLiter, UnitWeightMilliliter:
		return true
	default:
		return false
	}
}

// MarshalJSON for UnitWeightType
func (s UnitWeightType) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid unit type: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for UnitWeightType
func (s *UnitWeightType) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "KILOGRAM":
		*s = UnitWeightKilogram
	case "GRAM":
		*s = UnitWeightGram
	case "LITER":
		*s = UnitWeightLiter
	case "MILLILITER":
		*s = UnitWeightMilliliter
	default:
		return fmt.Errorf("invalid unit type: %s", str)
	}
	return nil
}
