package worker

// Source stores the code which will be transformed into another format.
type Source interface {
	String() string
}
