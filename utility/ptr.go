package utility

// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}
