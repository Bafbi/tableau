package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bafbi/tableau/internal/adapters/repo"
	"github.com/bafbi/tableau/internal/domain"
	"github.com/bafbi/tableau/internal/domain/query"
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
		
		slog.Debug("Creating task", "title", title)
		if err := r.Create(task); err != nil {
			return err
		}
		fmt.Printf("Created task %d: %s\n", task.ID, task.Title)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list [query]",
	Short: "List all tasks (e.g. 'status:todo priority:high')",
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
		
		var filter query.Filter
		if len(args) > 0 {
			filter = query.Parse(strings.Join(args, " "))
			slog.Debug("Filtering tasks", "query", args, "parsed", filter)
		} else {
			slog.Debug("Listing all tasks")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTitle\tStatus\tPriority")
		for _, t := range tasks {
			if len(args) > 0 && !query.Matches(t, filter) {
				continue
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, t.Title, t.Status, t.Priority)
		}
		w.Flush()
		return nil
	},
}

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a task in your editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		r := repo.NewFSRepository(cwd)
		
		task, err := r.Get(id)
		if err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		slog.Debug("Opening editor", "editor", editor, "file", task.FilePath)
		c := exec.Command(editor, task.FilePath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
}
