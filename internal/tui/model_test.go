package tui

import (
	"testing"

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

	if m.priorityIndex != 0 {
		t.Errorf("Expected priorityIndex 0, got %d", m.priorityIndex)
	}

	if m.titleInput.Placeholder != "Task title..." {
		t.Errorf("Expected title placeholder 'Task title...', got '%s'", m.titleInput.Placeholder)
	}

	if m.notesInput.Placeholder != "Add notes (optional)..." {
		t.Errorf("Expected notes placeholder 'Add notes (optional)...', got '%s'", m.notesInput.Placeholder)
	}
}
