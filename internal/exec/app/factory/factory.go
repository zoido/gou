package factory

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/lmittmann/tint"
	"github.com/zoido/gou/internal/dispatch"
	"github.com/zoido/gou/internal/exec/app/config"
)

// InputOpener is function that opens the open opens the configured input for reading.
//
// We do not return the reader directly to have dependency graph construction and app logic
// separated.
type InputOpener func() (io.ReadCloser, error)

// Factory is a struct for creating application components. It's a sole provider of the component
// instances.
type Factory struct {
	cfg config.Config

	// InputOpener returns function that opens configured input stream for reading.
	InputOpener func() InputOpener

	// Dispatcher returns action dispatcher.
	Dispatcher func() dispatch.Dispatcher

	// LogHandler returns configured [slog.Handler].
	SlogHandler func() slog.Handler
}

// New returns a new [Factory] using the provided [cfg].
func New(cfg config.Config) *Factory {
	f := &Factory{cfg: cfg}
	f.InputOpener = sync.OnceValue(f.createInputOpener)
	f.Dispatcher = sync.OnceValue(f.createDispatcher)
	f.SlogHandler = sync.OnceValue(f.createSlogHandler)
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

func (*Factory) createDispatcher() dispatch.Dispatcher {
	return make(dispatch.Dispatcher)
}

func (f *Factory) createSlogHandler() slog.Handler {
	level := slog.LevelInfo
	if f.cfg.Debug {
		level = slog.LevelDebug
	}
	return tint.NewHandler(os.Stderr, &tint.Options{Level: level})
}
