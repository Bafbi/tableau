package domain

import "time"

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
	CreatedAt   time.Time `toml:"created_at"`
	UpdatedAt   time.Time `toml:"updated_at"`
	Description string    `toml:"-"` // Content of the markdown file
	FilePath    string    `toml:"-"` // Path to the file
}

type TaskRepository interface {
	Init() error
	Create(task *Task) error
	List() ([]Task, error)
	Get(id int) (*Task, error)
	Update(task *Task) error
}
