package tui

import (
	"fmt"

	"github.com/bafbi/tableau/internal/domain"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Column struct {
	Title   string
	Status  domain.Status
	Tasks   []domain.Task
	Cursor  int
	Focused bool
	Width   int
	Height  int
}

func NewColumn(title string, status domain.Status, width, height int) Column {
	return Column{
		Title:  title,
		Status: status,
		Width:  width,
		Height: height,
	}
}

func (c *Column) SetTasks(tasks []domain.Task) {
	c.Tasks = tasks
	// Ensure cursor is valid
	if c.Cursor >= len(c.Tasks) && len(c.Tasks) > 0 {
		c.Cursor = len(c.Tasks) - 1
	}
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg) (Column, tea.Cmd) {
	// Handle local navigation if focused
	if !c.Focused {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if c.Cursor < len(c.Tasks)-1 {
				c.Cursor++
			}
		case "k", "up":
			if c.Cursor > 0 {
				c.Cursor--
			}
		}
	}
	return c, nil
}

func (c Column) View() string {
	var style lipgloss.Style
	if c.Focused {
		style = FocusedColumnStyle
	} else {
		style = ColumnStyle
	}

	// Render Title
	out := TitleStyle.Render(fmt.Sprintf("%s (%d)", c.Title, len(c.Tasks))) + "\n"

	// Render Tasks
	for i, t := range c.Tasks {
		prefix := "  "
		if t.Blocked {
			prefix = "🔒"
		}
		
		taskStr := fmt.Sprintf("#%d %s", t.ID, t.Title)
		if i == c.Cursor && c.Focused {
			out += SelectedTaskStyle.Render("> " + prefix + taskStr) + "\n"
		} else {
			out += TaskStyle.Render(prefix + taskStr) + "\n"
		}
	}

	// Fill remaining height
	// This is a bit naive, but works for MVP
	
	return style.Render(out)
}
