package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Project ProjectConfig `toml:"project"`
	Columns ColumnsConfig `toml:"columns"`
	Style   StyleConfig   `toml:"style"`
	Git     GitConfig     `toml:"git"`
}

type ProjectConfig struct {
	Name string `toml:"name"`
}

type ColumnsConfig struct {
	Todo   string `toml:"todo"`
	Doing  string `toml:"doing"`
	Review string `toml:"review"`
	Done   string `toml:"done"`
}

type StyleConfig struct {
	BorderColor   string `toml:"border_color"`
	SelectedColor string `toml:"selected_color"`
}

type GitConfig struct {
	AutoBranch   bool   `toml:"auto_branch"`
	BranchPrefix string `toml:"branch_prefix"`
}

func Default() Config {
	return Config{
		Project: ProjectConfig{Name: "Tableau"},
		Columns: ColumnsConfig{
			Todo:   "To Do",
			Doing:  "In Progress",
			Review: "Review",
			Done:   "Done",
		},
		Style: StyleConfig{
			BorderColor:   "#7D56F4",
			SelectedColor: "#ff00ff",
		},
		Git: GitConfig{
			AutoBranch:   true,
			BranchPrefix: "feat/",
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
