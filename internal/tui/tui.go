package tui

import (
	"fmt"
	"strings"

	"done/internal/task"

	tea "charm.land/bubbletea/v2"
)

func Run() error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}
	tasks, err := store.Load()
	if err != nil {
		return err
	}

	final, err := tea.NewProgram(model{
		store:    store,
		tasks:    tasks,
		finished: make(map[int]bool),
	}).Run()
	if err != nil {
		return err
	}
	if m, ok := final.(model); ok {
		var summary []string
		if m.added > 0 {
			label := "tasks"
			if m.added == 1 {
				label = "task"
			}
			summary = append(summary, fmt.Sprintf("added %d %s", m.added, label))
		}
		count := len(m.finished)
		if count > 0 {
			label := "tasks"
			if count == 1 {
				label = "task"
			}
			summary = append(summary, fmt.Sprintf("finished %d %s", count, label))
		}
		if len(summary) == 0 {
			fmt.Println("session: no tasks added or finished.")
		} else {
			fmt.Printf("session: %s.\n", strings.Join(summary, ", "))
		}
	}
	return nil
}
