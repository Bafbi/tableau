package main

import (
	"fmt"
	"os"

	"github.com/bafbi/tableau/internal/adapters/repo"
	"github.com/bafbi/tableau/internal/ui/tui"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Open the Kanban board",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		r := repo.NewFSRepository(cwd)
		
		// Ensure repo is initialized
		if err := r.Init(); err != nil {
			return err
		}

		p := tea.NewProgram(tui.NewModel(r), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(boardCmd)
}
