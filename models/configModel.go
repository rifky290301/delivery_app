package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// *-------------------------------------------------
// * NullString is a wrapper around sql.NullString for JSON unmarshalling
// *-------------------------------------------------

type NullString struct {
	sql.NullString
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	} else {
		ns.Valid = false
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// *-------------------------------------------------
// * NullTime is a wrapper around sql.NullTime to handle custom JSON marshalling
// *-------------------------------------------------
type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Time = *t
		nt.Valid = true
	} else {
		nt.Valid = false
	}
	return nil
}
