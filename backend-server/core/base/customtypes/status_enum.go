package customtypes

import (
	"fmt"
	"strings"
)

type StatusEnum int

const (
	_ StatusEnum = iota // Skip zero value to represent nil
	Enable
	Disable
	Draft
)

func (s StatusEnum) IsValid() bool {
	return s == Enable || s == Disable || s == Draft
}

func (s StatusEnum) String() string {
	switch s {
	case Enable:
		return "enable"
	case Disable:
		return "disable"
	case Draft:
		return "draft"
	default:
		return "unknown"
	}
}

// MarshalJSON for StatusEnum
func (s StatusEnum) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "unknown" {
		return nil, fmt.Errorf("invalid StatusEnum: %d", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for StatusEnum
func (s *StatusEnum) UnmarshalJSON(data []byte) error {
	str := strings.ToLower(string(data))
	switch str {
	case `"enable"`:
		*s = Enable
	case `"disable"`:
		*s = Disable
	case `"draft"`:
		*s = Draft
	default:
		return fmt.Errorf("invalid StatusEnum: %s", str)
	}
	return nil
}
