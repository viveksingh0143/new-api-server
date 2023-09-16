package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StockStatus string

const (
	STOCK_DISPATCHING StockStatus = "STOCK-DISPATCHING"
	STOCK_IN          StockStatus = "STOCK-IN"
	STOCK_OUT         StockStatus = "STOCK-OUT"
	STOCK_REJECT      StockStatus = "STOCK-REJECT"
	STOCK_UNKNOWN     StockStatus = "STOCK_UNKNOWN"
)

func GetAllStockStatus() []StockStatus {
	return []StockStatus{STOCK_DISPATCHING, STOCK_IN, STOCK_OUT, STOCK_REJECT}
}

func (s StockStatus) IsValid() bool {
	for _, validType := range GetAllStockStatus() {
		if s == validType {
			return true
		}
	}
	return false
}

func (s StockStatus) String() string {
	if s.IsValid() {
		return string(s)
	}
	return string(STOCK_UNKNOWN)
}

func ValidateStockStatus(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(StockStatus)
	if !ok {
		return false
	}
	return value.IsValid()
}

// MarshalJSON for StockStatus
func (s StockStatus) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid stock status: %s", s)
	}
	return []byte(`"` + s.String() + `"`), nil
}

// UnmarshalJSON for StockStatus
func (s *StockStatus) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))

	switch str {
	case strings.ToUpper(string(STOCK_DISPATCHING)):
		*s = STOCK_DISPATCHING
	case strings.ToUpper(string(STOCK_IN)):
		*s = STOCK_IN
	case strings.ToUpper(string(STOCK_OUT)):
		*s = STOCK_OUT
	case strings.ToUpper(string(STOCK_REJECT)):
		*s = STOCK_REJECT
	default:
		*s = STOCK_UNKNOWN
	}
	return nil
}
