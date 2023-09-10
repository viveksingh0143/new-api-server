package customtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type LoginViaEnum string

const (
	StaffID  LoginViaEnum = "STAFF_ID"
	Email    LoginViaEnum = "EMAIL"
	Username LoginViaEnum = "USERNAME"
)

func (s LoginViaEnum) String() string {
	switch s {
	case StaffID:
		return "STAFF_ID"
	case Email:
		return "EMAIL"
	case Username:
		return "USERNAME"
	default:
		return "UNKNOWN"
	}
}

func (e *LoginViaEnum) ViaEmail() bool {
	return *e == Email
}

func (e *LoginViaEnum) ViaUsername() bool {
	return *e == Username
}

func (e *LoginViaEnum) ViaStaffID() bool {
	return *e == StaffID
}

func ValidateLoginViaEnum(fl validator.FieldLevel) bool {
	rawValue, ok := fl.Field().Interface().(LoginViaEnum)
	if !ok {
		return false
	}
	value := LoginViaEnum(strings.ToUpper(string(rawValue)))

	switch value {
	case StaffID, Email, Username:
		return true
	default:
		return false
	}
}

// MarshalJSON for LoginViaEnum
func (s LoginViaEnum) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "UNKNOWN" {
		return nil, fmt.Errorf("invalid LoginViaEnum: %s", s)
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON for LoginViaEnum
func (s *LoginViaEnum) UnmarshalJSON(data []byte) error {
	str := strings.ToUpper(strings.Trim(string(data), `"`))
	switch str {
	case "STAFF_ID":
		*s = StaffID
	case "EMAIL":
		*s = Email
	case "USERNAME":
		*s = Username
	default:
		return fmt.Errorf("invalid LoginViaEnum: %s", str)
	}
	return nil
}
