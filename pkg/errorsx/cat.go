package errorsx

type Category int16

const (
	CategoryDuplicate         = iota // Duplicate indicates a specific type of conflict caused by existing data.
	CategoryNotExists                // Resource not found / missing
	CategoryUnauthorized             // Authentication failed / invalid credentials
	CategoryValidation               // Input or business rule validation failed
	CategoryFormat                   // Data format / parsing error
	CategoryTimeout                  // Operation timed out
	CategoryDependency               // External dependency failure (DB, API, etc.)
	CategoryDependencyTimeout        // External dependency timeout
	CategoryUnavailable              // Service temporarily unavailable
	CategoryForbidden                // Operation forbidden by business rules
)
