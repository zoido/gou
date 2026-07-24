package ui

import (
	"context"
	"errors"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/zoido/gou/internal/domain/finding"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// ErrNoSelection is returned when the user quits the picker without selecting an item.
var ErrNoSelection = errors.New("no finding selected")

type item struct {
	v *finding.Finding
}

func (i item) Title() string       { return i.v.URL() }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.v.URL() }

type List struct {
	list list.Model
	pick *finding.Finding
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeys(msg)
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m List) handleKeys(kp tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch kp.String() {
	case "q":
		return m, tea.Quit
	case "enter":
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.pick = i.v
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(kp)
	return m, cmd
}

func (m List) View() tea.View {
	v := tea.NewView(docStyle.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func Pick(ctx context.Context, findings []finding.Finding) (finding.Finding, error) {
	items := make([]list.Item, len(findings))
	for i, f := range findings {
		items[i] = item{v: &f}
	}
	l := &List{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}

	p := tea.NewProgram(l, tea.WithContext(ctx))
	finalModel, err := p.Run()
	if err != nil {
		return finding.Finding{}, err
	}

	m, ok := finalModel.(List)
	if !ok || m.pick == nil {
		return finding.Finding{}, ErrNoSelection
	}
	return *m.pick, nil
}
