package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/storage"
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

	if m.confirmPurge {
		t.Error("Expected confirmPurge to be false")
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

func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func TestPurgeCompleted(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Done task", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	if err := tl.Add("Pending task", models.Medium, ""); err != nil {
		t.Fatal(err)
	}
	tl.Todos[0].Completed = true

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// 'c' opens the purge confirmation modal.
	upd, _ := m.Update(keyMsg("c"))
	purgeModel, ok := upd.(model)
	if !ok {
		t.Fatalf("Expected model after Update, got %T", upd)
	}
	if !purgeModel.confirmPurge {
		t.Fatal("Expected confirmPurge to be true after pressing c")
	}
	if !strings.Contains(purgeModel.View(), "Purge all completed todos?") {
		t.Error("Expected purge confirmation text in the rendered view")
	}

	// 'n' cancels without removing anything.
	upd, _ = purgeModel.Update(keyMsg("n"))
	purgeModel, _ = upd.(model)
	if purgeModel.confirmPurge {
		t.Error("Expected confirmPurge to be false after pressing n")
	}
	if purgeModel.todoList.Count() != 2 {
		t.Errorf("Expected 2 todos after cancel, got %d", purgeModel.todoList.Count())
	}

	// 'c' then 'y' purges completed todos.
	upd, _ = purgeModel.Update(keyMsg("c"))
	purgeModel, _ = upd.(model)
	upd, _ = purgeModel.Update(keyMsg("y"))
	purgeModel, _ = upd.(model)
	if purgeModel.confirmPurge {
		t.Error("Expected confirmPurge to be false after confirming")
	}
	if purgeModel.todoList.Count() != 1 {
		t.Errorf("Expected 1 todo after purge, got %d", purgeModel.todoList.Count())
	}
	if purgeModel.todoList.Todos[0].Title != "Pending task" {
		t.Errorf("Expected 'Pending task' to remain, got '%s'", purgeModel.todoList.Todos[0].Title)
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

func TestSaveTodosCmd_SnapshotsOnCall(t *testing.T) {
	tl := models.NewTodoList()
	for _, task := range []struct {
		title    string
		priority models.Priority
		notes    string
	}{
		{"Buy groceries", models.Low, ""},
		{"Walk the dog", models.Medium, ""},
		{"Read", models.High, "chapter 3"},
	} {
		if err := tl.Add(task.title, task.priority, task.notes); err != nil {
			t.Fatal(err)
		}
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))
	cmd := m.saveTodosCmd()

	// Mutate the live list after the save command was created.
	if err := tl.Toggle(1); err != nil {
		t.Fatal(err)
	}
	if err := tl.Add("Should not appear", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	if err := tl.Delete(2); err != nil {
		t.Fatal(err)
	}

	msg := cmd()
	if msg == nil {
		t.Fatal("Expected a msgTodoSaved back from the save command")
	}
	if savedMsg, ok := msg.(msgTodoSaved); !ok {
		t.Fatalf("Expected msgTodoSaved, got %T", msg)
	} else if savedMsg.err != nil {
		t.Fatalf("Save failed: %v", savedMsg.err)
	}

	saved, err := storage.LoadTodos(m.filename)
	if err != nil {
		t.Fatalf("Loading saved file: %v", err)
	}
	if len(saved) != 3 {
		t.Fatalf("Snapshot should have 3 todos, got %d (live list was mutated after save cmd creation)", len(saved))
	}
	if saved[0].Completed {
		t.Error("Snapshot should reflect state at save command creation, not later toggles")
	}
}

func TestSaveTodosCmd_ConcurrentSaves(t *testing.T) {
	const todoCount = 32
	const iterations = 300

	tl := models.NewTodoList()
	for i := 0; i < todoCount; i++ {
		if err := tl.Add(fmt.Sprintf("todo %d", i), models.Low, ""); err != nil {
			t.Fatal(err)
		}
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	jobs := make(chan tea.Cmd)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cmd := range jobs {
				cmd()
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		if err := tl.Toggle(1 + i%todoCount); err != nil {
			t.Fatal(err)
		}
		jobs <- m.saveTodosCmd()
	}
	close(jobs)
	wg.Wait()
}
