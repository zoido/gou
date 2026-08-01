package dispatch

import (
	"os"

	osc52 "github.com/aymanbagabas/go-osc52/v2"
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
	case "y":
		_, err := osc52.New(f.URL()).WriteTo(os.Stderr)
		return false, err
	}

	return doQuit, nil
}
