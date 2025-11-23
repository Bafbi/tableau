package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err := root.Execute()
	return buf.String(), err
}

func TestCLIIntegration(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "tableau-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Override TABLEAU_DIR env var for the test process
	// Note: This affects the global state, so parallel tests might be an issue if we ran them.
	// But for this sequence it's fine.
	// However, the repo adapter reads os.Getenv("TABLEAU_DIR") inside NewFSRepository.
	// We need to make sure we set it before running commands.
	_ = os.Setenv("TABLEAU_DIR", ".tableau_test")
	
	// We need to change the CWD to the tmpDir so that NewFSRepository(cwd) works as expected
	// relative to the "project root".
	// Actually, NewFSRepository takes rootDir.
	// The commands use os.Getwd().
	// So we should change CWD to tmpDir.
	
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// 1. Init
	out, err := executeCommand(rootCmd, "init")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if !strings.Contains(out, "Initialized tableau project") {
		t.Errorf("init output unexpected: %s", out)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".tableau_test")); os.IsNotExist(err) {
		t.Error(".tableau_test directory not created")
	}

	// 2. New Task
	out, err = executeCommand(rootCmd, "new", "Integration Test Task")
	if err != nil {
		t.Fatalf("new failed: %v", err)
	}
	if !strings.Contains(out, "Created task 1") {
		t.Errorf("new output unexpected: %s", out)
	}
	taskFile := filepath.Join(tmpDir, ".tableau_test", "tasks", "1-integration-test-task.md")
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		t.Error("Task file not created")
	}

	// 3. List
	out, err = executeCommand(rootCmd, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "Integration Test Task") {
		t.Errorf("list output missing task: %s", out)
	}

	// 4. Block
	out, err = executeCommand(rootCmd, "block", "1")
	if err != nil {
		t.Fatalf("block failed: %v", err)
	}
	if !strings.Contains(out, "Task 1 blocked") {
		t.Errorf("block output unexpected: %s", out)
	}

	// 5. Comment
	out, err = executeCommand(rootCmd, "comment", "1", "Test comment")
	if err != nil {
		t.Fatalf("comment failed: %v", err)
	}
	if !strings.Contains(out, "Comment added") {
		t.Errorf("comment output unexpected: %s", out)
	}
}
