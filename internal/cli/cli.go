package cli

import (
	"fmt"
	"strings"

	"done/internal/task"
	"done/internal/tui"
)

const usage = `done - a terminal task finisher

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
	case "add":
		return add(args[1:])
	case "ls":
		return list()
	case "finish", "clear":
		return fmt.Errorf("%s is not implemented yet", args[0])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func add(args []string) error {
	title := strings.TrimSpace(strings.Join(args, " "))
	if title == "" {
		return fmt.Errorf("provide a task title\n\n%s", usage)
	}

	store, err := task.NewStore()
	if err != nil {
		return err
	}
	added, err := store.Add(title)
	if err != nil {
		return err
	}

	fmt.Printf("Added #%d: %s\n", added.ID, added.Title)
	return nil
}

func list() error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}
	tasks, err := store.Load()
	if err != nil {
		return err
	}

	open := 0
	for _, item := range tasks {
		if item.Completed {
			continue
		}
		fmt.Printf("%d. %s\n", item.ID, item.Title)
		open++
	}
	if open == 0 {
		fmt.Println("No open tasks.")
	}
	return nil
}
