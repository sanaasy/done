package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Store struct {
	Path string
}

func NewStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("find config directory: %w", err)
	}
	return &Store{Path: filepath.Join(configDir, "done", "tasks.json")}, nil
}

func (s *Store) Load() ([]Task, error) {
	file, err := os.Open(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open task list: %w", err)
	}
	defer file.Close()

	var tasks []Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		if errors.Is(err, io.EOF) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to decode task list %q: %w", s.Path, err)
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

func (s *Store) Save(tasks []Task) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return fmt.Errorf("failed to create task config directory: %w", err)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode task list: %w", err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(filepath.Dir(s.Path), ".tasks-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp task file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to set task file permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write task list: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close task list: %w", err)
	}
	if err := os.Rename(tmpPath, s.Path); err != nil {
		return fmt.Errorf("failed to save task list: %w", err)
	}
	return nil
}

func (s *Store) Add(title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("task title cannot be empty")
	}

	tasks, err := s.Load()
	if err != nil {
		return Task{}, err
	}

	nextID := 1
	for _, existing := range tasks {
		if existing.ID >= nextID {
			nextID = existing.ID + 1
		}
	}

	added := Task{ID: nextID, Title: title, CreatedAt: time.Now()}
	tasks = append(tasks, added)
	if err := s.Save(tasks); err != nil {
		return Task{}, err
	}
	return added, nil
}

func (s *Store) SetCompleted(id int, completed bool) (Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return Task{}, err
	}

	for i := range tasks {
		if tasks[i].ID != id {
			continue
		}
		if tasks[i].Completed != completed {
			tasks[i].Completed = completed
			if completed {
				tasks[i].CompletedAt = time.Now()
			} else {
				tasks[i].CompletedAt = time.Time{}
			}
			if err := s.Save(tasks); err != nil {
				return Task{}, err
			}
		}
		return tasks[i], nil
	}
	return Task{}, fmt.Errorf("task #%d not found", id)
}

func (s *Store) Clear() (int, error) {
	tasks, err := s.Load()
	if err != nil {
		return 0, err
	}

	remaining := make([]Task, 0, len(tasks))
	removed := 0
	for _, item := range tasks {
		if item.Completed {
			removed++
			continue
		}
		remaining = append(remaining, item)
	}
	if removed == 0 {
		return 0, nil
	}
	if err := s.Save(remaining); err != nil {
		return 0, err
	}
	return removed, nil
}
