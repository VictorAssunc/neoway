package pointer

// New creates a new pointer to v value.
func New[T any](v T) *T {
	return &v
}
