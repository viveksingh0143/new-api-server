package customtypes

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

const TimeFormat = "02-01-2006 03:04:05 PM"

type NullTime struct {
	sql.NullTime
}

func GetNullTime(value *time.Time) NullTime {
	if value != nil {
		return NewValidNullTime(*value)
	} else {
		return NewInvalidNullTime()
	}
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
	return json.Marshal(n.Time.Format(TimeFormat))
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullTime.
func (n *NullTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if s == "null" {
		n.Valid = false
		return nil
	}

	parsedTime, err := time.Parse(TimeFormat, s)
	if err != nil {
		log.Printf("%+v\n", err)
		n.Valid = false
		return err
	}

	n.Time = parsedTime
	n.Valid = true
	return nil
}
