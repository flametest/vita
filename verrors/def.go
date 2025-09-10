package verrors

var (
	ErrTypeAssertion = New(1001, "type assertion failed")
)

var predefinedErrors = []*Error{
	ErrTypeAssertion,
}
