package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#ff00ff")
	ColorGray      = lipgloss.Color("#626262")
	ColorLightGray = lipgloss.Color("#A8A8A8")

	// Styles
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	ColumnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGray).
			Height(20).
			Width(30)

	FocusedColumnStyle = ColumnStyle.
				BorderForeground(ColorPrimary)

	TaskStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ColorGray)

	SelectedTaskStyle = TaskStyle.
				BorderForeground(ColorSecondary).
				Foreground(ColorSecondary)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			PaddingBottom(1)
)
