package query

import (
	"testing"

	"github.com/bafbi/tableau/internal/domain"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected Filter
	}{
		{
			input: "status:todo",
			expected: Filter{
				Status: ptrStatus(domain.StatusTodo),
			},
		},
		{
			input: "priority:high",
			expected: Filter{
				Priority: ptrPriority(domain.PriorityHigh),
			},
		},
		{
			input: "label:bug label:ui",
			expected: Filter{
				Labels: []string{"bug", "ui"},
			},
		},
		{
			input: "status:done priority:low label:feature",
			expected: Filter{
				Status:   ptrStatus(domain.StatusDone),
				Priority: ptrPriority(domain.PriorityLow),
				Labels:   []string{"feature"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Parse(tt.input)
			if tt.expected.Status != nil && (got.Status == nil || *got.Status != *tt.expected.Status) {
				t.Errorf("Status: expected %v, got %v", tt.expected.Status, got.Status)
			}
			if tt.expected.Priority != nil && (got.Priority == nil || *got.Priority != *tt.expected.Priority) {
				t.Errorf("Priority: expected %v, got %v", tt.expected.Priority, got.Priority)
			}
			if len(tt.expected.Labels) > 0 {
				if len(got.Labels) != len(tt.expected.Labels) {
					t.Errorf("Labels: expected %v, got %v", tt.expected.Labels, got.Labels)
				}
			}
		})
	}
}

func TestMatches(t *testing.T) {
	task := domain.Task{
		Status:   domain.StatusTodo,
		Priority: domain.PriorityHigh,
		Labels:   []string{"bug", "backend"},
	}

	tests := []struct {
		name   string
		filter Filter
		match  bool
	}{
		{
			name:   "Match Status",
			filter: Filter{Status: ptrStatus(domain.StatusTodo)},
			match:  true,
		},
		{
			name:   "Mismatch Status",
			filter: Filter{Status: ptrStatus(domain.StatusDone)},
			match:  false,
		},
		{
			name:   "Match Priority",
			filter: Filter{Priority: ptrPriority(domain.PriorityHigh)},
			match:  true,
		},
		{
			name:   "Match Label",
			filter: Filter{Labels: []string{"bug"}},
			match:  true,
		},
		{
			name:   "Mismatch Label",
			filter: Filter{Labels: []string{"ui"}},
			match:  false,
		},
		{
			name:   "Match All",
			filter: Filter{Status: ptrStatus(domain.StatusTodo), Priority: ptrPriority(domain.PriorityHigh)},
			match:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(task, tt.filter); got != tt.match {
				t.Errorf("Matches() = %v, want %v", got, tt.match)
			}
		})
	}
}

func ptrStatus(s domain.Status) *domain.Status {
	return &s
}

func ptrPriority(p domain.Priority) *domain.Priority {
	return &p
}
