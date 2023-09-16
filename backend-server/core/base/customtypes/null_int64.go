package customtypes

import (
	"database/sql"
	"encoding/json"
	"log"
)

type NullInt64 struct {
	sql.NullInt64
}

// Initialize a valid NullInt64
func NewValidNullInt64(value int64) NullInt64 {
	return NullInt64{
		sql.NullInt64{Int64: value, Valid: true},
	}
}

// Initialize an invalid NullInt64
func NewInvalidNullInt64() NullInt64 {
	return NullInt64{
		sql.NullInt64{Valid: false},
	}
}

// MarshalJSON for custom JSON marshaling
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}

// UnmarshalJSON for custom JSON unmarshaling
func (n *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if i != nil {
		n.Int64 = *i
		n.Valid = true
	} else {
		n.Valid = false
	}
	return nil
}
