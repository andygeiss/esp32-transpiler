package worker

// Mapping specifies the api logic to apply transformation to a specific identifier.
type Mapping interface {
	Apply(ident string) string
}
