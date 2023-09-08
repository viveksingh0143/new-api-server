package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

// Initialize a valid NullTime
func NewValidNullTime(value time.Time) NullTime {
	return NullTime{
		sql.NullTime{Time: value, Valid: true},
	}
}

// Initialize an invalid NullTime
func NewInvalidNullTime() NullTime {
	return NullTime{
		sql.NullTime{Valid: false},
	}
}

// MarshalJSON implements the json.Marshaler interface for NullTime.
func (n *NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Time.Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullTime.
func (n *NullTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "null" {
		n.Valid = false
		return nil
	}

	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Time = parsedTime
	n.Valid = true
	return nil
}
