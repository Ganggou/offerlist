package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Jsonb is used to represent jsonb data structure in postgres
type Jsonb map[string]interface{}

// Value returns value of the Jsonb data
func (a Jsonb) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan ...
func (a *Jsonb) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
