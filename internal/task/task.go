// Package task defines the task data model used by the CLI and TUI.
package task

import "time"

// Task is one item in the user's task list.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}
