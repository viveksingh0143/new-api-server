package customtypes

import (
	"database/sql"
	"encoding/json"
	"log"
)

type NullFloat64 struct {
	sql.NullFloat64
}

// Initialize a valid NullFloat64
func NewValidNullFloat64(value float64) NullFloat64 {
	return NullFloat64{
		sql.NullFloat64{Float64: value, Valid: true},
	}
}

// Initialize an invalid NullFloat64
func NewInvalidNullFloat64() NullFloat64 {
	return NullFloat64{
		sql.NullFloat64{Valid: false},
	}
}

// MarshalJSON for custom JSON marshaling
func (n *NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}

// UnmarshalJSON for custom JSON unmarshaling
func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	var i *float64
	if err := json.Unmarshal(data, &i); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if i != nil {
		n.Float64 = *i
		n.Valid = true
	} else {
		n.Valid = false
	}
	return nil
}
