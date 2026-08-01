package dispatch

import (
	"github.com/rkoesters/xdg"

	"github.com/zoido/gou/internal/domain/finding"
)

type Dispatcher struct{}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(key string, f finding.Finding) (bool, error) {
	doQuit := false

	switch key {
	case "enter":
		return false, xdg.Open(f.URL())
	}

	return doQuit, nil
}
