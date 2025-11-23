package git

import (
	"log/slog"
	"os/exec"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CurrentBranch() (string, error) {
	slog.Debug("Running git rev-parse --abbrev-ref HEAD")
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) Checkout(branch string) error {
	slog.Debug("Running git checkout", "branch", branch)
	return exec.Command("git", "checkout", branch).Run()
}

func (c *Client) CreateBranch(branch string) error {
	slog.Debug("Running git checkout -b", "branch", branch)
	return exec.Command("git", "checkout", "-b", branch).Run()
}

func (c *Client) BranchExists(branch string) bool {
	slog.Debug("Running git rev-parse --verify", "branch", branch)
	err := exec.Command("git", "rev-parse", "--verify", branch).Run()
	return err == nil
}

func (c *Client) GetUserName() (string, error) {
	slog.Debug("Running git config user.name")
	out, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return "Unknown", nil
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "Unknown", nil
	}
	return name, nil
}
