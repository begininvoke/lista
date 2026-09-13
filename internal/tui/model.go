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
