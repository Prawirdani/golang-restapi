package request

// RequestBody is an interface for request body that needs to be validated and sanitized when binding to struct on handler function
// It is recommended to implement this interface on every request body model struct
// Should use this alongside with BindValidate function from binder.go
type RequestBody interface {
	// Validate is a method to validate request body to ensure all required fields are provided and match the constraints
	Validate() error
	// Sanitize is a method to sanitize request body to ensure all fields are in the correct format and clean
	Sanitize() error
}
