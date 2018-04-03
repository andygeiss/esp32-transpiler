package transpiler

// Transpiler specifies the behaviour of taking Go source code and transforming it into another language.
type Transpiler interface {
	Transpile() error
}
