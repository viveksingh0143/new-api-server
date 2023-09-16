package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StockLevel string

const (
	EMPTY_STOCK   StockLevel = "EMPTY"
	PARTIAL_STOCK StockLevel = "PARTIAL"
	FULL_STOCK    StockLevel = "FULL"
)

func GetAllStockLevels() []StockLevel {
	return []StockLevel{EMPTY_STOCK, PARTIAL_STOCK, FULL_STOCK}
}

func (s StockLevel) IsValid() bool {
	return s == EMPTY_STOCK || s == PARTIAL_STOCK || s == FULL_STOCK
}

func (s StockLevel) IsEmpty() bool {
	return s == EMPTY_STOCK
}
func (s StockLevel) IsPartial() bool {
	return s == PARTIAL_STOCK
}
func (s StockLevel) IsFull() bool {
	return s == FULL_STOCK
}

func (s StockLevel) String() string {
	switch s {
	case EMPTY_STOCK:
		return "EMPTY"
	case PARTIAL_STOCK:
		return "PARTIAL"
	case FULL_STOCK:
		return "FULL"
	default:
		return "UNKNOWN"
	}
}

func ValidateStockLevel(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(StockLevel)
	if !ok {
		return false
	}
	value := StockLevel(strings.ToUpper(string(rawValue)))

	switch value {
	case EMPTY_STOCK, PARTIAL_STOCK, FULL_STOCK:
		return true
	default:
		return false
	}
}

// MarshalJSON for StockLevel
func (s StockLevel) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid container type: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for StockLevel
func (s *StockLevel) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "EMPTY":
		*s = EMPTY_STOCK
	case "PARTIAL":
		*s = PARTIAL_STOCK
	case "FULL":
		*s = FULL_STOCK
	default:
		return fmt.Errorf("invalid container type: %s", str)
	}
	return nil
}
