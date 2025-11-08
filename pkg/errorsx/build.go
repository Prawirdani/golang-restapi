package errorsx

var (
	Duplicate         = build(CategoryDuplicate)
	NotExists         = build(CategoryNotExists)
	Unauthorized      = build(CategoryUnauthorized)
	Validation        = build(CategoryValidation)
	Format            = build(CategoryFormat)
	Timeout           = build(CategoryTimeout)
	Dependency        = build(CategoryDependency)
	DependencyTimeout = build(CategoryDependencyTimeout)
	Unavailable       = build(CategoryUnavailable)
	Forbidden         = build(CategoryForbidden)
)

func build(cat Category) func(msg string) *categorizedError {
	return func(msg string) *categorizedError {
		return &categorizedError{
			category: cat,
			Message:  msg,
		}
	}
}
