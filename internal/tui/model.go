package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

type model struct{}

func (model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyPressMsg); ok {
		return m, tea.Quit
	}
	return m, nil
}

func (model) View() tea.View {
	return tea.NewView(titleStyle.Render("done") + "\n\ntask list starter\n\nPress any key to quit.")
}
