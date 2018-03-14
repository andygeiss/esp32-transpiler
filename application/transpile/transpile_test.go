package transpile_test

import (
	"bytes"
	. "github.com/andygeiss/assert"
	"github.com/andygeiss/esp32-transpiler/application/transpile"
	"github.com/andygeiss/esp32-transpiler/business/worker"
	"io"
	"testing"
)

type mockupWorker struct {
	in  io.Reader
	out io.Writer
}

func (w *mockupWorker) Prepare(source []worker.Source) (chan worker.Source, error) {
	out := make(chan worker.Source)
	return out, nil
}
func (w *mockupWorker) Start() error {
	return nil
}
func (w *mockupWorker) Transform(source chan worker.Source) (chan worker.Target, error) {
	out := make(chan worker.Target)
	return out, nil
}

func TestTranspileErrorShouldBeNil(t *testing.T) {
	var in, out bytes.Buffer
	worker := &mockupWorker{&in, &out}
	trans := transpile.NewTranspiler(worker)
	Assert(t, trans.Transpile(), IsNil())
}
