package repo

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/adrg/frontmatter"
	"github.com/bafbi/tableau/internal/domain"
)

func ParseTask(r io.Reader) (*domain.Task, error) {
	var task domain.Task
	body, err := frontmatter.Parse(r, &task)
	if err != nil {
		return nil, err
	}
	task.Description = string(body)
	return &task, nil
}

func MarshalTask(task *domain.Task) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("+++\n")
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(task); err != nil {
		return nil, err
	}
	buf.WriteString("+++\n\n")
	buf.WriteString(task.Description)
	return buf.Bytes(), nil
}
