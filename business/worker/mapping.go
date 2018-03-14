package worker

// Mapping specifies the business logic to apply transformation to a specific identifier.
type Mapping interface {
	Apply(ident string) string
	Read() error
}
