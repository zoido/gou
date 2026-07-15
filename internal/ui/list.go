package ui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/zoido/gou/internal/domain/finding"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	v finding.Finding
}

func (i item) Title() string       { return i.v.URL() }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.v.URL() }

type List struct {
	list list.Model
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m List) View() tea.View {
	v := tea.NewView(docStyle.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func NewList(findings []finding.Finding) List {
	items := make([]list.Item, len(findings))
	for i, f := range findings {
		items[i] = item{v: f}
	}
	return List{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
}
