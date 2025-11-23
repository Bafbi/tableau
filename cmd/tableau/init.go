package main

import (
	"fmt"
	"os"

	"github.com/bafbi/tableau/internal/adapters/repo"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new tableau project",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		r := repo.NewFSRepository(cwd)
		if err := r.Init(); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Initialized tableau project in %s\n", r.DirName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
