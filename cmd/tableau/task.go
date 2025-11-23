package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bafbi/tableau/internal/adapters/repo"
	"github.com/bafbi/tableau/internal/domain"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		r := repo.NewFSRepository(cwd)
		
		title := strings.Join(args, " ")
		task := &domain.Task{
			Title:    title,
			Status:   domain.StatusTodo,
			Priority: domain.PriorityMedium,
		}
		
		if err := r.Create(task); err != nil {
			return err
		}
		fmt.Printf("Created task %d: %s\n", task.ID, task.Title)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		r := repo.NewFSRepository(cwd)
		tasks, err := r.List()
		if err != nil {
			return err
		}
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTitle\tStatus\tPriority")
		for _, t := range tasks {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, t.Title, t.Status, t.Priority)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
}
