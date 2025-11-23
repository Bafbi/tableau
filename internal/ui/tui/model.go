package tui

import (
	"github.com/bafbi/tableau/internal/domain"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Repo       domain.TaskRepository
	Columns    []Column
	FocusedCol int
	Quitting   bool
	Err        error
}

func NewModel(repo domain.TaskRepository) Model {
	return Model{
		Repo: repo,
		Columns: []Column{
			NewColumn("To Do", domain.StatusTodo, 30, 20),
			NewColumn("In Progress", domain.StatusDoing, 30, 20),
			NewColumn("Done", domain.StatusDone, 30, 20),
		},
		FocusedCol: 0,
	}
}

type tasksLoadedMsg []domain.Task
type errMsg error

func (m Model) Init() tea.Cmd {
	return m.loadTasks
}

func (m Model) loadTasks() tea.Msg {
	tasks, err := m.Repo.List()
	if err != nil {
		return errMsg(err)
	}
	return tasksLoadedMsg(tasks)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "h", "left":
			if m.FocusedCol > 0 {
				m.Columns[m.FocusedCol].Focused = false
				m.FocusedCol--
				m.Columns[m.FocusedCol].Focused = true
			}
		case "l", "right":
			if m.FocusedCol < len(m.Columns)-1 {
				m.Columns[m.FocusedCol].Focused = false
				m.FocusedCol++
				m.Columns[m.FocusedCol].Focused = true
			}
		case "H": // Move task left
			return m.moveTask(-1)
		case "L": // Move task right
			return m.moveTask(1)
		}
	case tasksLoadedMsg:
		// Distribute tasks to columns
		// Reset columns first
		for i := range m.Columns {
			m.Columns[i].Tasks = []domain.Task{}
		}
		
		for _, t := range msg {
			switch t.Status {
			case domain.StatusTodo:
				m.Columns[0].Tasks = append(m.Columns[0].Tasks, t)
			case domain.StatusDoing:
				m.Columns[1].Tasks = append(m.Columns[1].Tasks, t)
			case domain.StatusDone:
				m.Columns[2].Tasks = append(m.Columns[2].Tasks, t)
			}
		}
		
		// Ensure focus is set
		m.Columns[m.FocusedCol].Focused = true
		
		// Refresh cursors
		for i := range m.Columns {
			m.Columns[i].SetTasks(m.Columns[i].Tasks)
		}
		
	case errMsg:
		m.Err = msg
		return m, tea.Quit
	}

	// Update focused column
	var cmd tea.Cmd
	m.Columns[m.FocusedCol], cmd = m.Columns[m.FocusedCol].Update(msg)
	return m, cmd
}

func (m Model) moveTask(direction int) (tea.Model, tea.Cmd) {
	col := &m.Columns[m.FocusedCol]
	if len(col.Tasks) == 0 {
		return m, nil
	}
	
	task := col.Tasks[col.Cursor]
	newColIdx := m.FocusedCol + direction
	
	if newColIdx < 0 || newColIdx >= len(m.Columns) {
		return m, nil
	}
	
	// Update task status
	task.Status = m.Columns[newColIdx].Status
	
	// Save to repo
	if err := m.Repo.Update(&task); err != nil {
		m.Err = err
		return m, nil // Or show error
	}
	
	// Reload tasks to refresh view
	return m, m.loadTasks
}

func (m Model) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error()
	}
	
	var cols []string
	for _, c := range m.Columns {
		cols = append(cols, c.View())
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}
