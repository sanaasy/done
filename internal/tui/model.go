package tui

import (
	"fmt"
	"strings"
	"time"

	"done/internal/task"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a1ff"))
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#777784"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#d6b8ff"))
	doneStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#79cba0")).Strikethrough(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f28b82"))
)

type model struct {
	store    *task.Store
	tasks    []task.Task
	cursor   int
	adding   bool
	input    string
	notice   string
	added    int
	finished map[int]bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		if m.adding {
			return m.updateInput(key)
		}
		m.notice = ""
		switch key.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "space":
			m.toggle()
		case "a":
			m.adding = true
			m.input = ""
		}
	}
	return m, nil
}

func (m model) updateInput(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.adding = false
		m.input = ""
	case "enter":
		added, err := m.store.Add(m.input)
		if err != nil {
			m.notice = err.Error()
			return m, nil
		}
		m.tasks = append(m.tasks, added)
		m.added++
		m.cursor = len(m.tasks) - 1
		m.adding = false
		m.input = ""
	case "backspace":
		text := []rune(m.input)
		if len(text) > 0 {
			m.input = string(text[:len(text)-1])
		}
	case "ctrl+c":
		return m, tea.Quit
	default:
		m.input += key.Key().Text
	}
	return m, nil
}

func (m *model) toggle() {
	if len(m.tasks) == 0 || m.cursor >= len(m.tasks) {
		return
	}
	item, err := m.store.SetCompleted(m.tasks[m.cursor].ID, !m.tasks[m.cursor].Completed)
	if err != nil {
		m.notice = err.Error()
		return
	}
	m.tasks[m.cursor] = item
	if item.Completed {
		m.finished[item.ID] = true
	} else {
		delete(m.finished, item.ID)
	}
}

func (m model) View() tea.View {
	var out strings.Builder
	done := 0
	for _, item := range m.tasks {
		if item.Completed {
			done++
		}
	}
	fmt.Fprintf(&out, "%s %s\n\n", titleStyle.Render("tasks"), mutedStyle.Render(fmt.Sprintf("(%d/%d)", done, len(m.tasks))))

	if len(m.tasks) == 0 {
		out.WriteString(mutedStyle.Render("press a to add a task."))
		out.WriteString("\n")
	} else {
		for i, item := range m.tasks {
			branch := "├─"
			if i == len(m.tasks)-1 {
				branch = "└─"
			}
			marker, style := "○", mutedStyle
			if item.Completed {
				marker, style = "✓", doneStyle
			} else if i == m.cursor {
				marker, style = "◐", selectedStyle
			}
			pointer := "  "
			if i == m.cursor {
				pointer = "› "
			}
			fmt.Fprintf(&out, "%s%s %s %s\n", pointer, mutedStyle.Render(branch), style.Render(marker), style.Render(item.Title))
		}
	}

	finishedToday, today := todayProgress(m.tasks, time.Now())
	out.WriteString("\n")
	fmt.Fprintf(&out, "%s %d/%d\n", mutedStyle.Render("Today"), finishedToday, today)
	out.WriteString(progress(finishedToday, today, 28))
	out.WriteString("\n\n")
	if m.adding {
		out.WriteString(titleStyle.Render("add task: ") + m.input + "▏  " + mutedStyle.Render("enter add · esc cancel"))
	} else {
		out.WriteString(mutedStyle.Render("↑/↓ move   space toggle   a add   q quit"))
	}
	if m.notice != "" {
		out.WriteString("\n" + errorStyle.Render(m.notice))
	}
	return tea.NewView(out.String())
}

func todayProgress(tasks []task.Task, now time.Time) (int, int) {
	done, total := 0, 0
	for _, item := range tasks {
		createdToday := sameDay(item.CreatedAt, now)
		completedToday := item.Completed && !item.CompletedAt.IsZero() && sameDay(item.CompletedAt, now)
		if createdToday || completedToday {
			total++
		}
		if completedToday {
			done++
		}
	}
	return done, total
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.In(b.Location()).Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func progress(done, total, width int) string {
	filled := 0
	if total > 0 {
		filled = done * width / total
	}
	return selectedStyle.Render(strings.Repeat("━", filled)) + mutedStyle.Render(strings.Repeat("─", width-filled))
}
