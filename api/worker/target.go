package worker

// Target stores the result of the transformation process.
type Target interface {
	String() string
}
