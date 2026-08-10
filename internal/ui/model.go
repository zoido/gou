package ui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/zoido/gou/internal/model"
)

const extraHeight = 4 // title=2, help=2

type item struct {
	v *model.Finding
}

func (i item) Title() string       { return i.v.URL() }
func (item) Description() string   { return "" }
func (i item) FilterValue() string { return i.v.URL() }

func newListModel(findings []model.Finding) list.Model {
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

	return l
}

type wrapper struct {
	list     list.Model
	findings []model.Finding
	dispatch DispatchFn
}

func (wrapper) Init() tea.Cmd {
	return nil
}

func (m wrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeys(msg)
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(len(m.findings)+extraHeight, msg.Height))
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m wrapper) handleKeys(kp tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	i := m.list.Index()
	handled, doQuit, err := m.dispatch(kp.String(), m.findings[i])
	if doQuit {
		cmd = tea.Quit
	}
	if err != nil {
		m.list.NewStatusMessage(err.Error())
	}
	if handled {
		return m, cmd
	}

	m.list, cmd = m.list.Update(kp)
	return m, cmd
}

func (m wrapper) View() tea.View {
	v := tea.NewView(m.list.View())
	v.AltScreen = false
	return v
}
