// Package errorsx provides a type-safe, category-based error system.
//
// It defines error Categories to represent general types of problems
// within the application. Categories are abstract and business-oriented,
// keeping the domain layer decoupled from transport concerns like HTTP.
//
// Common sentinel errors can be exposed as exported variables or functions
// for convenience, while the underlying constructor remains private to
// ensure consistent creation.
package errorsx

type CategorizedError interface {
	error
	Category() Category
}

type categorizedError struct {
	Message  string
	category Category
}

func (e *categorizedError) Error() string {
	return e.Message
}

func (e *categorizedError) Category() Category {
	return e.category
}

func New(message string, category Category) *categorizedError {
	return &categorizedError{
		Message:  message,
		category: category,
	}
}
