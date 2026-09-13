package cmd

import (
	"testing"

	"github.com/kwame-Owusu/lista/internal/models"
)

func sampleTodos() []models.Todo {
	return []models.Todo{
		{ID: 1, Title: "Buy milk", Notes: "", Completed: false, Priority: models.High},
		{ID: 2, Title: "Write report", Notes: "quarterly", Completed: true, Priority: models.Low},
		{ID: 3, Title: "Call plumber", Notes: "", Completed: false, Priority: models.Medium},
	}
}

func TestFilterTodos(t *testing.T) {
	tests := []struct {
		name   string
		filter listFilter
		want   []int
	}{
		{
			name:   "no filter keeps all",
			filter: listFilter{},
			want:   []int{1, 2, 3},
		},
		{
			name:   "pending only",
			filter: listFilter{pending: true},
			want:   []int{1, 3},
		},
		{
			name:   "done only",
			filter: listFilter{done: true},
			want:   []int{2},
		},
		{
			name:   "priority high",
			filter: listFilter{priority: priorityPtr(models.High)},
			want:   []int{1},
		},
		{
			name:   "priority low",
			filter: listFilter{priority: priorityPtr(models.Low)},
			want:   []int{2},
		},
		{
			name:   "search matches title",
			filter: listFilter{search: "plumber"},
			want:   []int{3},
		},
		{
			name:   "search matches notes",
			filter: listFilter{search: "quarterly"},
			want:   []int{2},
		},
		{
			name:   "search is case-insensitive",
			filter: listFilter{search: "MILK"},
			want:   []int{1},
		},
		{
			name:   "search no match",
			filter: listFilter{search: "nothing"},
			want:   []int{},
		},
		{
			name:   "pending plus priority",
			filter: listFilter{pending: true, priority: priorityPtr(models.High)},
			want:   []int{1},
		},
		{
			name:   "search plus done",
			filter: listFilter{done: true, search: "quarterly"},
			want:   []int{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterTodos(sampleTodos(), tt.filter)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d todos, want %d", len(got), len(tt.want))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("todo[%d].ID = %d, want %d", i, got[i].ID, id)
				}
			}
		})
	}
}

func priorityPtr(p models.Priority) *models.Priority {
	return &p
}
