package git

import (
	"os/exec"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) Checkout(branch string) error {
	return exec.Command("git", "checkout", branch).Run()
}

func (c *Client) CreateBranch(branch string) error {
	return exec.Command("git", "checkout", "-b", branch).Run()
}

func (c *Client) BranchExists(branch string) bool {
	err := exec.Command("git", "rev-parse", "--verify", branch).Run()
	return err == nil
}
