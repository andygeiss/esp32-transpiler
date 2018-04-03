package worker

// Worker specifies the api logic of transforming a source code format into another target format.
type Worker interface {
	Start() error
}
