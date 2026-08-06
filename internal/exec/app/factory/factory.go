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
	"github.com/zoido/gou/internal/ui"
)

// InputOpener is function that opens the open opens the configured input for reading.
//
// We do not return the reader directly to have dependency graph construction and app logic
// separated.
type InputOpener func() (io.ReadCloser, error)

// Factory is a struct for creating application components. It's a sole provider of the component
// instances.
type Factory struct {
	cfg *config.Config

	dispatcher func() dispatch.Dispatcher

	// InputOpener returns function that opens configured input stream for reading.
	InputOpener func() InputOpener

	// LogHandler returns configured [slog.Handler].
	SlogHandler func() slog.Handler

	// Program returns instance of the program.
	Program func() *ui.Program
}

// New returns a new [Factory] using the provided [cfg].
func New(cfg *config.Config) *Factory {
	f := &Factory{cfg: cfg}

	f.dispatcher = sync.OnceValue(f.createDispatcher)

	f.InputOpener = sync.OnceValue(f.createInputOpener)
	f.SlogHandler = sync.OnceValue(f.createSlogHandler)
	f.Program = sync.OnceValue(f.createProgram)

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

func (f *Factory) createProgram() *ui.Program {
	d := f.dispatcher()
	return ui.NewProgram(d.Dispatch)
}

func (f *Factory) createSlogHandler() slog.Handler {
	level := slog.LevelInfo
	if f.cfg.Debug {
		level = slog.LevelDebug
	}
	return tint.NewTextHandler(os.Stderr, &tint.Options{Level: level})
}
