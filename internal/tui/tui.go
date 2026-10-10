package tui

import tea "charm.land/bubbletea/v2"

func Run() error {
	_, err := tea.NewProgram(model{}).Run()
	return err
}
