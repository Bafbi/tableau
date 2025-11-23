package repo

import (
	"bytes"
	"testing"
	"time"

	"github.com/bafbi/tableau/internal/domain"
)

func TestMarshalAndParseTask(t *testing.T) {
	original := &domain.Task{
		ID:          1,
		Title:       "Test Task",
		Status:      domain.StatusDoing,
		Priority:    domain.PriorityHigh,
		Labels:      []string{"bug"},
		Assignee:    "me",
		CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		Description: "# Context\nThis is a test task.",
		Comments: []domain.Comment{
			{
				Author:    "tester",
				Text:      "A comment",
				CreatedAt: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Marshal
	data, err := MarshalTask(original)
	if err != nil {
		t.Fatalf("MarshalTask failed: %v", err)
	}

	// Parse
	parsed, err := ParseTask(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseTask failed: %v", err)
	}

	// Verify
	if parsed.ID != original.ID {
		t.Errorf("ID: expected %d, got %d", original.ID, parsed.ID)
	}
	if parsed.Title != original.Title {
		t.Errorf("Title: expected %s, got %s", original.Title, parsed.Title)
	}
	if parsed.Status != original.Status {
		t.Errorf("Status: expected %s, got %s", original.Status, parsed.Status)
	}
	if len(parsed.Comments) != 1 {
		t.Errorf("Comments: expected 1, got %d", len(parsed.Comments))
	}
	if parsed.Comments[0].Text != original.Comments[0].Text {
		t.Errorf("Comment Text: expected %s, got %s", original.Comments[0].Text, parsed.Comments[0].Text)
	}
	if parsed.Description != original.Description {
		t.Errorf("Description: expected %q, got %q", original.Description, parsed.Description)
	}
}
