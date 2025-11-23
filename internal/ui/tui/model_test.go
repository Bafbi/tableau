package tui

import (
	"testing"

	"github.com/bafbi/tableau/internal/config"
	"github.com/bafbi/tableau/internal/domain"
	"github.com/charmbracelet/bubbletea"
)

// MockRepo implements TaskRepository for testing
type MockRepo struct {
	tasks []domain.Task
}

func (m *MockRepo) Init() error                        { return nil }
func (m *MockRepo) LoadConfig() (config.Config, error) { return config.Default(), nil }
func (m *MockRepo) Create(task *domain.Task) error     { return nil }
func (m *MockRepo) List() ([]domain.Task, error)       { return m.tasks, nil }
func (m *MockRepo) Get(id int) (*domain.Task, error)   { return nil, nil }
func (m *MockRepo) Update(task *domain.Task) error     { return nil }

func TestModelUpdate(t *testing.T) {
	repo := &MockRepo{
		tasks: []domain.Task{
			{ID: 1, Title: "Task 1", Status: domain.StatusTodo},
			{ID: 2, Title: "Task 2", Status: domain.StatusDoing},
		},
	}

	model := NewModel(repo)

	// Simulate Init (loading tasks)
	// We manually trigger the tasksLoadedMsg because Init returns a Cmd
	tasks, _ := repo.List()
	model.Update(tasksLoadedMsg(tasks))

	// Verify initial state
	if len(model.Columns[0].Tasks) != 1 {
		t.Errorf("Expected 1 task in Todo, got %d", len(model.Columns[0].Tasks))
	}
	if model.FocusedCol != 0 {
		t.Errorf("Expected focused col 0, got %d", model.FocusedCol)
	}

	// Test Navigation (Right)
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight, Runes: []rune{'l'}})
	m := newModel.(Model)
	if m.FocusedCol != 1 {
		t.Errorf("Expected focused col 1 after moving right, got %d", m.FocusedCol)
	}

	// Test Navigation (Left)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft, Runes: []rune{'h'}})
	m = newModel.(Model)
	if m.FocusedCol != 0 {
		t.Errorf("Expected focused col 0 after moving left, got %d", m.FocusedCol)
	}

	// Test Detail View
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(Model)
	if m.State != DetailView {
		t.Errorf("Expected state DetailView after Enter, got %v", m.State)
	}

	// Test Escape from Detail View
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(Model)
	if m.State != BoardView {
		t.Errorf("Expected state BoardView after Esc, got %v", m.State)
	}
}
