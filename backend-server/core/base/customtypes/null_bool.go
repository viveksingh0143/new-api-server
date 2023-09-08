package customtypes

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

// Initialize a valid NullBool
func NewValidNullBool(value bool) NullBool {
	return NullBool{
		sql.NullBool{Bool: value, Valid: true},
	}
}

// Initialize an invalid NullBool
func NewInvalidNullBool() NullBool {
	return NullBool{
		sql.NullBool{Valid: false},
	}
}

// MarshalJSON implements json.Marshaler interface
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}
