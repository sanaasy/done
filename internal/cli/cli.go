// Package cli routes command-line arguments to the task store or interactive UI.
package cli

import (
	"errors"
	"fmt"

	"done/internal/tui"
)

const usage = `done — a tiny terminal task finisher

Usage:
  done                         Open the task list
  done add "Task description"  Add a task
  done ls                      List outstanding tasks
  done finish <id>             Mark a task complete
  done clear                   Remove completed tasks
`

func Run(args []string) error {
	if len(args) == 0 {
		return tui.Run()
	}

	switch args[0] {
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	case "add", "ls", "finish", "clear":
		return errors.New("in progress")
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}
