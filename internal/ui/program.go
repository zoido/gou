package ui

import (
	"context"
	"errors"

	tea "charm.land/bubbletea/v2"

	"github.com/zoido/gou/internal/model"
)

// DispatchFn performs action with the finding based on the pressed key.
type DispatchFn func(string, model.Finding) (handled, doQuit bool, err error)

// Program displays list of findings and calls [DispatchFn] on key presses.
type Program struct {
	dispatch DispatchFn
}

// NewProgram returns [Program] instance using [d] for dispatching actions.
func NewProgram(d DispatchFn) *Program {
	return &Program{dispatch: d}
}

// Run executes the [Program].
func (p *Program) Run(ctx context.Context, findings []model.Finding) error {
	modl := &wrapper{
		list:     newListModel(findings),
		dispatch: p.dispatch,
		findings: findings,
	}
	prg := tea.NewProgram(modl, tea.WithContext(ctx))

	_, err := prg.Run()
	if errors.Is(err, tea.ErrProgramKilled) || ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
