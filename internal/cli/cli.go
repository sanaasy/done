package cli

import (
	"fmt"
	"strconv"
	"strings"

	"done/internal/task"
	"done/internal/tui"
)

const usage = `done - a terminal task finisher

Usage:
  done                         open the task list
  done add "task description"  add a task
  done ls                      list outstanding tasks
  done check <id>              mark a task complete
  done clear                   remove completed tasks
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
	case "check":
		return check(args[1:])
	case "clear":
		return clear(args[1:])
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

	fmt.Printf("added #%d: %s\n", added.ID, added.Title)
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
		fmt.Println("no open tasks.")
	}
	return nil
}

func check(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: done check <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id < 1 {
		return fmt.Errorf("task ID must be a positive num")
	}
	store, err := task.NewStore()
	if err != nil {
		return err
	}
	item, err := store.SetCompleted(id, true)
	if err != nil {
		return err
	}
	fmt.Printf("finished #%d: %s\n", item.ID, item.Title)
	return nil
}

func clear(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: done clear")
	}
	store, err := task.NewStore()
	if err != nil {
		return err
	}
	removed, err := store.Clear()
	if err != nil {
		return err
	}
	if removed == 0 {
		fmt.Println("no tasks to clear.")
		return nil
	}
	label := "tasks"
	if removed == 1 {
		label = "task"
	}
	fmt.Printf("cleared %d %s.\n", removed, label)
	return nil
}
