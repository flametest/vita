package vtool

// Deprecated: use github.com/samber/lo.ToPtr instead.
func Ptr[T any](v T) *T {
	return &v
}
