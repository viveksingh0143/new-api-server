package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type NullStatusEnum struct {
	Status StatusEnum
	Valid  bool
}

func (nse NullStatusEnum) Value() (driver.Value, error) {
	if !nse.Valid {
		return nil, nil
	}
	return nse.Status, nil
}

// Initialize a valid NullStatusType
func NewValidNullStatusType(value StatusEnum) NullStatusEnum {
	return NullStatusEnum{
		Status: value,
		Valid:  true,
	}
}

// Initialize an invalid NullStatusType
func NewInvalidNullStatusType() NullStatusEnum {
	return NullStatusEnum{
		Valid: false,
	}
}

// MarshalJSON for custom JSON marshaling
func (n *NullStatusEnum) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Status.String())
}

// UnmarshalJSON for custom JSON unmarshaling
func (n *NullStatusEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	mappedStatus := MapStatusEnum(s)
	if mappedStatus == InvalidStatus {
		n.Valid = false
		return errors.New("invalid status value")
	}

	n.Status = mappedStatus
	n.Valid = true
	return nil
}
