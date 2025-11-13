package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Nullable represents an optional value of type T.
// Valid indicates whether Val is set (non-zero).
type Nullable[T comparable] struct {
	val   T
	valid bool
}

// New returns a Nullable wrapping value.
// If value equals the zero value of T, Valid is false.
func New[T comparable](value T) Nullable[T] {
	var zero T
	if value == zero {
		return Nullable[T]{val: zero, valid: false}
	}
	return Nullable[T]{val: value, valid: true}
}

// NotNull reports whether the Nullable holds a valid (non-zero) value.
func (n Nullable[T]) NotNull() bool {
	return n.valid
}

// Scan implements the [sql.Scanner] interface, called when scanning a row from the database.
func (n *Nullable[T]) Scan(value any) error {
	var zero T
	if value == nil {
		n.val, n.valid = zero, false
		return nil
	}

	v, ok := value.(T)
	if !ok {
		// Try marshaling
		if v, ok := value.(string); ok {
			if err := json.Unmarshal([]byte(v), &n.val); err != nil {
				return err
			}
			if n.val != zero {
				n.valid = true
			}
			return nil
		}
		return fmt.Errorf("cannot convert %T to %T", value, v)
	}

	n.val = v
	n.valid = true
	return nil
}

// MarshalJSON implements the [json.Marshaler] interface, called when marshaling to JSON.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.val)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface, called when unmarshaling from JSON.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.valid = false
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.val = v
	n.valid = true
	return nil
}

// Value implements [driver.Valuer]. It is called automatically when the value
// is written to the database (e.g., via Exec or Query arguments).
//
// This allows passing Nullable directly to database operations without an
// explicit getter.
//
// Do not use this method for value comparisons or assertions — use Get() to
// retrieve the underlying value instead.
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.valid {
		return nil, nil
	}
	return n.val, nil
}

// Get returns the wrapped value if valid, or the zero value of T otherwise.
func (n Nullable[T]) Get() T {
	return n.val
}

// Set sets the value of the nullable, if value is zero value of T, it will set the nullable to invalid.
// TODO: Should allow zero value, eg int with 0 should be a valid one, perhaps provide option to set the validity
func (n *Nullable[T]) Set(value T) {
	var zero T
	if value == zero {
		n.valid = false
		n.val = zero
		return
	}

	n.val = value
	n.valid = true
}
