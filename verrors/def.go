package verrors

var (
	ErrTypeAssertion  = New(1001, "type assertion failed")
	ErrOptimisticLock = New(1002, "optimistic lock conflict")
)

var predefinedErrors = []*Error{
	ErrTypeAssertion,
	ErrOptimisticLock,
}
