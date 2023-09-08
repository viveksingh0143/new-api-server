package customtypes

import (
	"database/sql/driver"
	"encoding/json"
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

// MarshalJSON implements the json.Marshaler interface for NullStatusEnum.
func (nse NullStatusEnum) MarshalJSON() ([]byte, error) {
	if !nse.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nse.Status)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullStatusEnum.
func (nse *NullStatusEnum) UnmarshalJSON(data []byte) error {
	var temp *StatusEnum

	// Try to unmarshal into the temporary pointer
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// If unmarshal was successful and the temporary pointer is not nil,
	// then set the status enum as valid
	if temp != nil {
		nse.Status = *temp
		nse.Valid = true
	} else {
		nse.Valid = false
	}

	return nil
}
