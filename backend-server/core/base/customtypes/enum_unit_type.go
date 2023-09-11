package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UnitType string

const (
	UnitKilogram   UnitType = "Kilogram"
	UnitGram       UnitType = "Gram"
	UnitLiter      UnitType = "Liter"
	UnitMilliliter UnitType = "Milliliter"
	UnitPiece      UnitType = "Piece"
)

func GetAllUnitTypes() []UnitType {
	return []UnitType{UnitKilogram, UnitGram, UnitLiter, UnitMilliliter, UnitPiece}
}

func (s UnitType) IsValid() bool {
	return s == UnitKilogram || s == UnitGram || s == UnitLiter || s == UnitMilliliter || s == UnitPiece
}

func (s UnitType) String() string {
	switch s {
	case UnitKilogram:
		return "Kilogram"
	case UnitGram:
		return "Gram"
	case UnitLiter:
		return "Liter"
	case UnitMilliliter:
		return "Milliliter"
	case UnitPiece:
		return "Piece"
	default:
		return "UNKNOWN"
	}
}

func ValidateUnitType(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(UnitType)
	if !ok {
		return false
	}
	value := UnitType(string(rawValue))

	switch value {
	case UnitKilogram, UnitGram, UnitLiter, UnitMilliliter, UnitPiece:
		return true
	default:
		return false
	}
}

// MarshalJSON for UnitType
func (s UnitType) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid unit type: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for UnitType
func (s *UnitType) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "KILOGRAM":
		*s = UnitKilogram
	case "GRAM":
		*s = UnitGram
	case "LITER":
		*s = UnitLiter
	case "MILLILITER":
		*s = UnitMilliliter
	case "PIECE":
		*s = UnitPiece
	default:
		return fmt.Errorf("invalid unit type: %s", str)
	}
	return nil
}
