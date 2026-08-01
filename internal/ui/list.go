package ui

import (
	"context"
	"errors"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/zoido/gou/internal/domain/finding"
)

// ErrNoSelection is returned when the user quits the picker without selecting an item.
var ErrNoSelection = errors.New("no finding selected")

const extraHeight = 4 // title=2, help=2

type item struct {
	v *finding.Finding
}

func (i item) Title() string       { return i.v.URL() }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.v.URL() }

type dispatchFn func(string, finding.Finding) (doQuit bool, err error)

func Serve(
	ctx context.Context,
	findings []finding.Finding,
	dispatch dispatchFn,
) error {
	items := make([]list.Item, len(findings))
	for i, f := range findings {
		items[i] = item{v: &f}
	}

	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)

	l := list.New(items, d, 0, 0)

	l.Title = "GoU"
	l.SetShowTitle(true)

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)

	p := tea.NewProgram(
		&Model{
			listModel: l,
			dispatch:  dispatch,
			findings:  findings,
		},
		tea.WithContext(ctx),
	)

	_, err := p.Run()
	if errors.Is(err, tea.ErrProgramKilled) {
		return ctx.Err()
	}
	if err != nil {
		return err
	}

	return nil
}

type Model struct {
	listModel list.Model
	findings  []finding.Finding
	dispatch  dispatchFn
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeys(msg)
	case tea.WindowSizeMsg:
		m.listModel.SetSize(msg.Width, min(len(m.findings)+extraHeight, msg.Height))
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)
	return m, cmd
}

func (m Model) handleKeys(kp tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if kp.String() == "q" {
		return m, tea.Quit
	}

	i := m.listModel.Index()
	doQuit, err := m.dispatch(kp.String(), m.findings[i])
	if doQuit {
		return m, tea.Quit
	}
	if err != nil {
		m.listModel.NewStatusMessage(err.Error())
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(kp)
	return m, cmd
}

func (m Model) View() tea.View {
	v := tea.NewView(m.listModel.View())
	v.AltScreen = false
	return v
}
