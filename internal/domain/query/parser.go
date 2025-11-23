package query

import (
	"strings"

	"github.com/bafbi/tableau/internal/domain"
)

type Filter struct {
	Status   *domain.Status
	Priority *domain.Priority
	Labels   []string
}

func Parse(q string) Filter {
	f := Filter{}
	parts := strings.Fields(q)
	for _, part := range parts {
		if strings.HasPrefix(part, "status:") {
			val := domain.Status(strings.TrimPrefix(part, "status:"))
			f.Status = &val
		} else if strings.HasPrefix(part, "priority:") {
			val := domain.Priority(strings.TrimPrefix(part, "priority:"))
			f.Priority = &val
		} else if strings.HasPrefix(part, "label:") {
			val := strings.TrimPrefix(part, "label:")
			f.Labels = append(f.Labels, val)
		}
	}
	return f
}

func Matches(t domain.Task, f Filter) bool {
	if f.Status != nil && t.Status != *f.Status {
		return false
	}
	if f.Priority != nil && t.Priority != *f.Priority {
		return false
	}
	for _, label := range f.Labels {
		found := false
		for _, tLabel := range t.Labels {
			if tLabel == label {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
