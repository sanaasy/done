// Package tui contains the interactive Bubble Tea interface.
package tui

import tea "charm.land/bubbletea/v2"

// Run starts the interactive task list.
func Run() error {
	_, err := tea.NewProgram(model{}).Run()
	return err
}
