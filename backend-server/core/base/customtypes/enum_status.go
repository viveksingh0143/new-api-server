package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StatusEnum int

const (
	_ StatusEnum = iota
	Enable
	Disable
	Draft
)

func (s StatusEnum) IsValid() bool {
	return s == Enable || s == Disable || s == Draft
}

func (s StatusEnum) IsEnable() bool {
	return s == Enable
}

func (s StatusEnum) IsDisable() bool {
	return s == Disable
}

func (s StatusEnum) IsDraft() bool {
	return s == Draft
}

func (s StatusEnum) String() string {
	switch s {
	case Enable:
		return "ENABLE"
	case Disable:
		return "DISABLE"
	case Draft:
		return "DRAFT"
	default:
		return "UNKNOWN"
	}
}

// MarshalJSON for StatusEnum
func (s StatusEnum) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid StatusEnum: %d", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for StatusEnum
func (s *StatusEnum) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "ENABLE":
		*s = Enable
	case "DISABLE":
		*s = Disable
	case "DRAFT":
		*s = Draft
	default:
		return fmt.Errorf("invalid StatusEnum: %s", str)
	}
	return nil
}

func ValidateStatusEnum(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(StatusEnum)
	if !ok {
		return false
	}
	return rawValue.IsValid()
}
