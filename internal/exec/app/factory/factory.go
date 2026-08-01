package factory

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/zoido/gou/internal/dispatch"
	"github.com/zoido/gou/internal/exec/app/config"
)

type InputOpener func() (io.ReadCloser, error)

// Factory is a struct for creating application components. It's a sole provider of the component
// instances.
type Factory struct {
	cfg config.Config

	// InputOpener returns function that opens configured input stream for reading.
	InputOpener func() InputOpener

	// Dispatcher returns action dispatcher.
	Dispatcher func() *dispatch.Dispatcher
}

// New returns a new [Factory] using the provided [cfg].
func New(cfg config.Config) *Factory {
	f := &Factory{cfg: cfg}
	f.InputOpener = sync.OnceValue(f.createInputOpener)
	f.Dispatcher = sync.OnceValue(f.createDispatcher)
	return f
}

func (f *Factory) createInputOpener() InputOpener {
	if f.cfg.File == "" || f.cfg.File == "-" {
		return func() (io.ReadCloser, error) {
			return os.Stdin, nil
		}
	}

	return func() (io.ReadCloser, error) {
		abs, err := filepath.Abs(f.cfg.File)
		if err != nil {
			return nil, err
		}
		return os.Open(abs)
	}
}

func (f *Factory) createDispatcher() *dispatch.Dispatcher {
	return dispatch.NewDispatcher()
}
