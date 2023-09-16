package customtypes

import (
	"database/sql"
	"encoding/json"
	"log"
)

type NullString struct {
	sql.NullString
}

func NewNullString(value string) NullString {
	return NullString{
		sql.NullString{String: value, Valid: value == ""},
	}
}

// Initialize a valid NullString
func NewValidNullString(value string) NullString {
	return NullString{
		sql.NullString{String: value, Valid: true},
	}
}

// Initialize an invalid NullString
func NewInvalidNullString() NullString {
	return NullString{
		sql.NullString{Valid: false},
	}
}

// MarshalJSON implements json.Marshaler interface
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
