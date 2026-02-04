package update

// TypedUpdate represents a typed update handle.
// This is a helper to satisfy the requirement of using update.New[Req, Res].
type TypedUpdate[Req any, Res any] struct {
	name string
}

// New creates a new TypedUpdate.
func New[Req any, Res any](name string) *TypedUpdate[Req, Res] {
	return &TypedUpdate[Req, Res]{name: name}
}

// Name returns the name of the update.
func (t *TypedUpdate[Req, Res]) Name() string {
	return t.name
}
