package dispatch

import (
	"fmt"
	"strings"

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
		if strings.HasPrefix(f.URL(), "gemini://") {
			return false, fmt.Errorf("gemini not supported")
		}
	}

	return doQuit, nil
}
