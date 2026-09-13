package tui

import (
	"cmp"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/models"
)

type tickMsg struct{}

const tickInterval = 30 * time.Second

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type formField int

const (
	fieldTitle formField = iota
	fieldPriority
	fieldNotes
)

const maxUndoEntries = 10

type undoEntry struct {
	snapshot models.TodoList
	cursor   int
}

type model struct {
	todoList      *models.TodoList
	cursor        int // which todo is selected
	width         int // terminal width
	height        int // terminal height
	filename      string
	err           error
	confirmDelete bool
	deleteID      int
	confirmPurge  bool
	showHelp      bool
	undoStack     []undoEntry

	// Form state
	addingTodo   bool
	focusedField formField
	titleInput   textinput.Model
	notesInput   textarea.Model
	priority     models.Priority

	// editing state
	editingTodo bool
	editingID   int
}

func NewModel(todoList *models.TodoList, filename string) model {
	// Title input
	ti := textinput.New()
	ti.Placeholder = "Task title..."
	ti.CharLimit = models.MaxTitleLength
	ti.Width = 50

	// Notes textarea
	ta := textarea.New()
	ta.Placeholder = "Add notes (optional)..."
	ta.CharLimit = models.MaxNotesLength
	ta.SetWidth(50)
	ta.SetHeight(5)
	ta.ShowLineNumbers = false

	return model{
		todoList:     todoList,
		cursor:       0,
		filename:     filename,
		titleInput:   ti,
		notesInput:   ta,
		priority:     models.Low,
		addingTodo:   false,
		focusedField: fieldTitle,
	}
}

// cloneTodoList returns a copy of the todo list so an undo snapshot can be
// restored without aliasing the live list. Todo is made of value fields, so
// copying the slice is sufficient for isolation.
func cloneTodoList(tl *models.TodoList) models.TodoList {
	return models.TodoList{
		Todos:  append([]models.Todo(nil), tl.Todos...),
		NextID: tl.NextID,
	}
}

// pushUndo records the current list state and cursor so the most recent
// toggle, delete, or edit can be reverted. The stack is bounded to
// maxUndoEntries, evicting the oldest snapshot when full.
func (m *model) pushUndo() {
	m.undoStack = append(m.undoStack, undoEntry{
		snapshot: cloneTodoList(m.todoList),
		cursor:   m.cursor,
	})
	if len(m.undoStack) > maxUndoEntries {
		m.undoStack = m.undoStack[1:]
	}
}

// performUndo restores the most recent undoable action and re-saves to disk.
// It returns nil when there is nothing to undo.
func (m *model) performUndo() tea.Cmd {
	if len(m.undoStack) == 0 {
		return nil
	}

	entry := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]

	m.todoList.Todos = append([]models.Todo(nil), entry.snapshot.Todos...)
	m.todoList.NextID = entry.snapshot.NextID
	m.cursor = entry.cursor
	return m.saveTodosCmd()
}

func (m model) Init() tea.Cmd {
	return tick()
}

// findTodoIndexByID returns the index of the todo with the given id using
// binary search. Assumes todos is sorted ascending by ID.
func findTodoIndexByID(todos []models.Todo, id int) int {
	idx, found := slices.BinarySearchFunc(todos, id, func(t models.Todo, id int) int {
		return cmp.Compare(t.ID, id)
	})
	if !found {
		return -1
	}
	return idx
}

var priorityOptions = []models.Priority{models.Low, models.Medium, models.High}

func priorityIndex(p models.Priority) int {
	for i, pri := range priorityOptions {
		if pri == p {
			return i
		}
	}
	return 0
}

func priorityAt(i int) models.Priority {
	if i < 0 || i >= len(priorityOptions) {
		return models.Low
	}
	return priorityOptions[i]
}
