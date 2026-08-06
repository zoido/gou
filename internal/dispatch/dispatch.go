package dispatch

import (
	"github.com/zoido/gou/internal/model"
)

// Dispatcher decides what action to perform based on the key press.
type Dispatcher map[string]Action

// Action represents single action done with th single [model.Finding]. Returned true indicates that
// application should end when done.
type Action func(model.Finding) (bool, error)

// Dispatch performs the action associated with the key press and returns its result.
func (d Dispatcher) Dispatch(key string, f model.Finding) (bool, bool, error) {
	a, ok := d[key]
	if !ok {
		return false, false, nil
	}
	doQuit, err := a(f)
	if err != nil {
		return true, false, err
	}
	return true, doQuit, nil
}
