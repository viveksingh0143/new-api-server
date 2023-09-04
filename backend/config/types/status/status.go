package status

type StatusType int

const (
	Disabled StatusType = iota
	Enabled
	Draft
)

func (s StatusType) String() string {
	return [...]string{"Disabled", "Enabled", "Draft"}[s]
}

func MapStatus(status string) StatusType {
	switch status {
	case "Disabled":
		return Disabled
	case "Enabled":
		return Enabled
	case "Draft":
		return Draft
	default:
		return Draft
	}
}
