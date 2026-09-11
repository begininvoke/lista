package tui

import (
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

	// Form state
	addingTodo    bool
	focusedField  formField
	titleInput    textinput.Model
	notesInput    textarea.Model
	priorityIndex int // 0=Low, 1=Medium, 2=High

	// editing state
	editingTodo bool
	editingID   int
}

func NewModel(todoList *models.TodoList, filename string) model {
	// Title input
	ti := textinput.New()
	ti.Placeholder = "Task title..."
	ti.CharLimit = 200
	ti.Width = 50

	// Notes textarea
	ta := textarea.New()
	ta.Placeholder = "Add notes (optional)..."
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(5)
	ta.ShowLineNumbers = false

	return model{
		todoList:      todoList,
		cursor:        0,
		filename:      filename,
		titleInput:    ti,
		notesInput:    ta,
		priorityIndex: 0, // Default to Low
		addingTodo:    false,
		focusedField:  fieldTitle,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func findTodoIndexByID(todos []models.Todo, id int) int {
	for i, t := range todos {
		if t.ID == id {
			return i
		}
	}
	return -1
}
