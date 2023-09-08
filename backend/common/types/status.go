package types

type StatusEnum int

const (
	DisabledStatus StatusEnum = iota
	EnabledStatus
	DraftStatus
	InvalidStatus
)

func (s StatusEnum) String() string {
	return [...]string{"Disabled", "Enabled", "Draft", "Invalid"}[s]
}

func MapStatusEnum(status string) StatusEnum {
	switch status {
	case "Disabled":
		return DisabledStatus
	case "Enabled":
		return EnabledStatus
	case "Draft":
		return DraftStatus
	default:
		return InvalidStatus
	}
}
