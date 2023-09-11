package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ProductType string

const (
	RAW_MATERIAL_TYPE        ProductType = "RAW Material"
	FINISHED_GOODS_TYPE      ProductType = "Finished Goods"
	SEMI_FINISHED_GOODS_TYPE ProductType = "Semi Finished Goods"
	UNKNOWN_TYPE             ProductType = "UNKNOWN"
)

func GetAllProductTypes() []ProductType {
	return []ProductType{RAW_MATERIAL_TYPE, FINISHED_GOODS_TYPE, SEMI_FINISHED_GOODS_TYPE}
}

func (s ProductType) IsValid() bool {
	for _, validType := range GetAllProductTypes() {
		if s == validType {
			return true
		}
	}
	return false
}

func (s ProductType) String() string {
	if s.IsValid() {
		return string(s)
	}
	return string(UNKNOWN_TYPE)
}

func ValidateProductType(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(ProductType)
	if !ok {
		return false
	}
	return value.IsValid()
}

// MarshalJSON for ProductType
func (s ProductType) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid product type: %s", s)
	}
	return []byte(`"` + s.String() + `"`), nil
}

// UnmarshalJSON for ProductType
func (s *ProductType) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))

	switch str {
	case strings.ToUpper(string(RAW_MATERIAL_TYPE)):
		*s = RAW_MATERIAL_TYPE
	case strings.ToUpper(string(FINISHED_GOODS_TYPE)):
		*s = FINISHED_GOODS_TYPE
	case strings.ToUpper(string(SEMI_FINISHED_GOODS_TYPE)):
		*s = SEMI_FINISHED_GOODS_TYPE
	default:
		*s = UNKNOWN_TYPE
	}
	return nil
}
