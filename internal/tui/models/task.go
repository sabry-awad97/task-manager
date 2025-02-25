package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	DueDate     time.Time     `json:"due_date"`
	Priority    PriorityLevel `json:"priority"`
	Completed   bool          `json:"completed"`
	CreatedAt   time.Time     `json:"created_at"`
}

type PriorityLevel int

const (
	Low PriorityLevel = iota
	Medium
	High
)

func NewTask(title string, description string, dueDate time.Time, priority PriorityLevel) Task {
	return Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Priority:    priority,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}

func (p PriorityLevel) String() string {
	return [...]string{"Low", "Medium", "High"}[p]
}

func (p PriorityLevel) Color() string {
	return [...]string{"#44B556", "#FFA500", "#FF0000"}[p]
}
