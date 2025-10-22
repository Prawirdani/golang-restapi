package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Nullable[T comparable] struct {
	Val   T
	Valid bool
}

// NewNullable creates a new nullable with the given value. if the value is zero, the nullable will be invalid.
func NewNullable[T comparable](value T) Nullable[T] {
	var zero T
	if value == zero {
		return Nullable[T]{Val: zero, Valid: false}
	}

	return Nullable[T]{Val: value, Valid: true}
}

// Scan implements the [Scanner] interface, called when scanning a row from the database.
func (n *Nullable[T]) Scan(value any) error {
	var zero T
	if value == nil {
		n.Val, n.Valid = zero, false
		return nil
	}

	v, ok := value.(T)
	if !ok {
		// Try marshaling
		if v, ok := value.(string); ok {
			if err := json.Unmarshal([]byte(v), &n.Val); err != nil {
				return err
			}
			if n.Val != zero {
				n.Valid = true
			}
			return nil
		}
		return fmt.Errorf("cannot convert %T to %T", value, v)
	}

	n.Val = v
	n.Valid = true
	return nil
}

// MarshalJSON implements the [json.Marshaler] interface, called when marshaling to JSON.
func (n *Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface, called when unmarshaling from JSON.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.Val = v
	n.Valid = true
	return nil
}

// Value implements the [driver.Valuer] interface. Called when inserting into the database eg exec function.
// This make when passing the struct to the exec function args, using it directly without calling additional getter function.
// When you need to comparing/asserting the value, use the nullable.Val field instead of the nullable itself.
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

// Set sets the value of the nullable, if the value is zero/default empty value, it will set the nullable to invalid.
func (n *Nullable[T]) Set(value T) {
	var zero T
	if value == zero {
		n.Valid = false
		n.Val = zero
		return
	}

	n.Val = value
	n.Valid = true
}
