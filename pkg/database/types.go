package database

// types.go contains Null types to be used for scanning into
// all types implement a custom MarshalJSON() function to remove valid from json

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"reflect"
	"time"
)

// NullString is a wrapper for sql.NullString for JSON methods to remove Valid
type NullString struct {
	sql.NullString
}

func ToNullString(value string, setValidOnEmpty bool) NullString {
	valid := true
	if value == "" && !setValidOnEmpty {
		valid = false
	}
	return NullString{sql.NullString{String: value, Valid: valid}}
}

// MarshalJSON converts the value to null or JSON string
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		s.Valid = false
		*ns = NullString{s}
	} else {
		s.Valid = true
		*ns = NullString{s}
	}

	return nil
}

// NullInt64 is a wrapper for sql.NullInt64 for JSON methods to remove Valid
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON converts the value to null or JSON int
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n NullInt64) String() string {
	return fmt.Sprintf("%d", n.Int64)
}

func ToNullInt(value int64, setValidOnEmpty bool) NullInt64 {
	valid := true
	if value == 0 && !setValidOnEmpty {
		valid = false
	}
	return NullInt64{sql.NullInt64{Int64: value, Valid: valid}}
}

// NullFloat64 is a wrapper for sql.NullFLoat64 for JSON methods to remove Valid
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON converts the value to null or JSON float
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

func (n NullFloat64) String() string {
	return fmt.Sprintf("%f", n.Float64)
}

// NullTime is a wrapper for mysql.NullTime for JSON methods to remove Valid
type NullTime struct {
	mysql.NullTime
}

func ToNullTime(value time.Time, setValidOnEmpty bool) NullTime {
	valid := true
	if value.IsZero() && !setValidOnEmpty {
		valid = false
	}
	return NullTime{mysql.NullTime{Time: value, Valid: valid}}
}

// MarshalJSON converts the value to null or JSON time
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}

	return json.Marshal(nil)
}

func (n NullTime) String() string {
	return n.Time.String()
}

// NullBool is a wrapper for sql.NullBool for JSON methods to remove Valid
type NullBool struct {
	sql.NullBool
}

// MarshalJSON converts the value to null or JSON bool
func (n *NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n NullBool) String() string {
	return fmt.Sprintf("%t", n.Bool)
}

// IntBool is a type wrapper for tinyint columns used as bools
type IntBool struct {
	Bool  bool
	Valid bool
}

// Scan converts 1 to true and 0 to false
func (n *IntBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}

	n.Valid = true

	switch t := value.(type) {
	case int64:
		if t == 0 {
			n.Bool = false
		}

		if t == 1 {
			n.Bool = true
		}
	default:
		return fmt.Errorf("unsupported Scan, for IntBool type %T into type %T", t, n.Bool)
	}

	return nil
}

// Value returns 0 for false and 1 for true
func (n IntBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	if n.Bool {
		return int64(1), nil
	}

	return int64(0), nil
}

// MarshalJSON converts the value to null or JSON bool
func (n IntBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON converts 1, 0, true, false to bool value from JSON or nil on null
func (n *IntBool) UnmarshalJSON(data []byte) error {
	if id, err := strconv.ParseBool(string(data)); err == nil {
		n.Valid = true
		n.Bool = id
		return nil
	}

	if string(data) == "null" {
		n.Valid = false
		n.Bool = false
		return nil
	}

	return errors.New("unmarshal invalid format for IntBool")
}

// String is the Stringer interface for IntBool
func (n IntBool) String() string {
	return fmt.Sprintf("%t", n.Bool)
}
