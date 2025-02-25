package storage

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/sabry-awad97/task-manager/internal/tui/models"
)

type JSONStore struct {
	filePath string
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{filePath: path}
}

func (s *JSONStore) Save(tasks []models.Task) error {
	data, err := json.MarshalIndent(tasks, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *JSONStore) Load() ([]models.Task, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Task{}, nil
		}
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
