package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/kwame-Owusu/lista/internal/models"
)

func TestNewModel(t *testing.T) {
	tl := models.NewTodoList()
	m := NewModel(tl, "test.json")

	if m.todoList != tl {
		t.Error("NewModel() did not store the todo list")
	}

	if m.filename != "test.json" {
		t.Errorf("Expected filename 'test.json', got '%s'", m.filename)
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor 0, got %d", m.cursor)
	}

	if m.addingTodo {
		t.Error("Expected addingTodo to be false")
	}

	if m.editingTodo {
		t.Error("Expected editingTodo to be false")
	}

	if m.confirmDelete {
		t.Error("Expected confirmDelete to be false")
	}

	if m.focusedField != fieldTitle {
		t.Errorf("Expected focusedField to be fieldTitle, got %v", m.focusedField)
	}

	if m.priority != models.Low {
		t.Errorf("Expected priority %v, got %v", models.Low, m.priority)
	}

	if m.titleInput.Placeholder != "Task title..." {
		t.Errorf("Expected title placeholder 'Task title...', got '%s'", m.titleInput.Placeholder)
	}

	if m.notesInput.Placeholder != "Add notes (optional)..." {
		t.Errorf("Expected notes placeholder 'Add notes (optional)...', got '%s'", m.notesInput.Placeholder)
	}
}

func TestPriorityMapping(t *testing.T) {
	for _, p := range priorityOptions {
		if got := priorityAt(priorityIndex(p)); got != p {
			t.Errorf("priorityAt(priorityIndex(%v)) = %v, want %v", p, got, p)
		}
	}
}

func TestCyclePriority(t *testing.T) {
	m := NewModel(models.NewTodoList(), "test.json")

	if m.priority != models.Low {
		t.Fatalf("Expected starting priority Low, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.Medium {
		t.Errorf("Expected Medium after cycling down, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.High {
		t.Errorf("Expected High after second cycle down, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.Low {
		t.Errorf("Expected wrap to Low, got %v", m.priority)
	}

	m.cyclePriority(true)
	if m.priority != models.High {
		t.Errorf("Expected High after cycling up from Low, got %v", m.priority)
	}
}

func TestTickRefresh(t *testing.T) {
	tl := models.NewTodoList()
	err := tl.Add("Test todo", models.Low, "")
	if err != nil {
		t.Fatal(err)
	}
	tl.Todos[0].CreatedAt = time.Now().Add(-2 * time.Second)

	m := NewModel(tl, "test.json")

	if space := strings.TrimSpace(m.View()); strings.Contains(space, "No todos yet") {
		t.Fatal("Expected rendered list, got 'No todos yet'")
	}

	updated, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Error("Expected Update to re-arm the tick command")
	}

	view := updated.View()
	if !strings.Contains(view, "added 2s ago") {
		t.Errorf("Expected fresh timestamp 'added 2s ago' in view, got:\n%s", view)
	}

	updatedForm, cmdForm := updated.Update(tickMsg{})
	if cmdForm == nil {
		t.Error("Expected heartbeat to stay armed while in a form")
	}
	if _, ok := updatedForm.(model); !ok {
		t.Errorf("Expected model back after tick, got %T", updatedForm)
	}
}
