package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/bafbi/tableau/internal/adapters/git"
	"github.com/bafbi/tableau/internal/adapters/repo"
	"github.com/bafbi/tableau/internal/domain"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [id]",
	Short: "Start working on a task (creates git branch)",
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

		g := git.NewClient()
		
		branchName := task.Branch
		if branchName == "" {
			slug := strings.ToLower(strings.ReplaceAll(task.Title, " ", "-"))
			branchName = fmt.Sprintf("feat/%d-%s", task.ID, slug)
		}

		if g.BranchExists(branchName) {
			slog.Debug("Branch exists, checking out", "branch", branchName)
			fmt.Printf("Switching to branch %s\n", branchName)
			if err := g.Checkout(branchName); err != nil {
				return err
			}
		} else {
			slog.Debug("Creating new branch", "branch", branchName)
			fmt.Printf("Creating branch %s\n", branchName)
			if err := g.CreateBranch(branchName); err != nil {
				return err
			}
		}

		// Update task
		task.Status = domain.StatusDoing
		task.Branch = branchName
		if err := r.Update(task); err != nil {
			return err
		}
		
		fmt.Printf("Task %d started\n", task.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
