package models

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	DueDate     time.Time
	Priority    PriorityLevel
	Completed   bool
}

type PriorityLevel int

const (
	Low PriorityLevel = iota
	Medium
	High
)

func (p PriorityLevel) String() string {
	return [...]string{"Low", "Medium", "High"}[p]
}
