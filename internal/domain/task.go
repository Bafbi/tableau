package domain

import (
	"time"

	"github.com/bafbi/tableau/internal/config"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

type Task struct {
	ID          int       `toml:"id"`
	Title       string    `toml:"title"`
	Status      Status    `toml:"status"`
	Priority    Priority  `toml:"priority"`
	Labels      []string  `toml:"labels"`
	Assignee    string    `toml:"assignee,omitempty"`
	Branch      string    `toml:"branch,omitempty"`
	Blocked     bool      `toml:"blocked"`
	Comments    []Comment `toml:"comments,omitempty"`
	CreatedAt   time.Time `toml:"created_at"`
	UpdatedAt   time.Time `toml:"updated_at"`
	Description string    `toml:"-"` // Content of the markdown file
	FilePath    string    `toml:"-"` // Path to the file
}

type Comment struct {
	Author    string    `toml:"author"`
	Text      string    `toml:"text"`
	CreatedAt time.Time `toml:"created_at"`
}

type TaskRepository interface {
	Init() error
	LoadConfig() (config.Config, error)
	Create(task *Task) error
	List() ([]Task, error)
	Get(id int) (*Task, error)
	Update(task *Task) error
}
