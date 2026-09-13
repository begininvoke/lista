package models

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxTitleLength = 200
	MaxNotesLength = 500
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("todo title cannot be empty")
	}
	if utf8.RuneCountInString(title) > MaxTitleLength {
		return fmt.Errorf("todo title exceeds %d characters", MaxTitleLength)
	}
	return nil
}

func validateNotes(notes string) error {
	if utf8.RuneCountInString(notes) > MaxNotesLength {
		return fmt.Errorf("todo notes exceed %d characters", MaxNotesLength)
	}
	return nil
}

func validatePriority(priority Priority) error {
	if !priority.IsValid() {
		return fmt.Errorf("invalid priority: %d", priority)
	}
	return nil
}

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes,omitempty"`
	Completed bool      `json:"completed"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Todo) TimeAgo() string {
	if t.CreatedAt.IsZero() {
		return ""
	}
	d := time.Since(t.CreatedAt)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("added %ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("added %dmins ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("added %dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("added %dd ago", int(d.Hours()/24))
	}
}

type TodoList struct {
	Todos  []Todo `json:"todos"`
	NextID int    `json:"nextID"`
}

func NewTodoList() *TodoList {
	return &TodoList{
		Todos:  make([]Todo, 0),
		NextID: 1,
	}
}

func (tl *TodoList) Add(title string, priority Priority, notes string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	if err := validatePriority(priority); err != nil {
		return err
	}
	if err := validateNotes(notes); err != nil {
		return err
	}
	todo := Todo{
		ID:        tl.NextID,
		Title:     title,
		Notes:     notes,
		Completed: false,
		Priority:  priority,
		CreatedAt: time.Now(),
	}
	tl.Todos = append(tl.Todos, todo)
	tl.NextID++
	return nil
}

// findIndex returns the index of the todo with the given id using binary
// search. Assumes tl.Todos is sorted ascending by ID.
func (tl *TodoList) findIndex(id int) int {
	idx, found := slices.BinarySearchFunc(tl.Todos, id, func(t Todo, id int) int {
		return cmp.Compare(t.ID, id)
	})
	if !found {
		return -1
	}
	return idx
}

func (tl *TodoList) Update(id int, title string, priority Priority, notes string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	if err := validatePriority(priority); err != nil {
		return err
	}
	if err := validateNotes(notes); err != nil {
		return err
	}

	i := tl.findIndex(id)
	if i < 0 {
		return errors.New("todo not found")
	}
	tl.Todos[i].Title = title
	tl.Todos[i].Priority = priority
	tl.Todos[i].Notes = notes
	return nil
}

func (tl *TodoList) Complete(id int) error {
	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	tl.Todos[i].Completed = true
	return nil
}

func (tl *TodoList) Uncomplete(id int) error {
	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	tl.Todos[i].Completed = false
	return nil
}

// List returns a copy of the todos so callers can read or mutate the
// snapshot without affecting the live list (e.g. sorting in CLI output).
func (tl *TodoList) List() []Todo {
	result := make([]Todo, len(tl.Todos))
	copy(result, tl.Todos)
	return result
}

func (tl *TodoList) GetByID(id int) (*Todo, error) {
	if len(tl.Todos) == 0 {
		return nil, fmt.Errorf("no entries in todo list")
	}

	i := tl.findIndex(id)
	if i < 0 {
		return nil, fmt.Errorf("todo item with id %d not found", id)
	}
	return &tl.Todos[i], nil
}

func (tl *TodoList) Delete(id int) error {
	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
	return nil
}

func (tl *TodoList) Edit(id int, title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}

	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	tl.Todos[i].Title = title
	return nil
}

func (tl *TodoList) AppendNotes(id int, notes string) error {
	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	newNotes := notes
	if tl.Todos[i].Notes != "" {
		newNotes = tl.Todos[i].Notes + " " + notes
	}
	if err := validateNotes(newNotes); err != nil {
		return err
	}
	tl.Todos[i].Notes = newNotes
	return nil
}

func (tl *TodoList) Count() int {
	return len(tl.Todos)
}

func (tl *TodoList) GetPending() []Todo {
	if tl.Count() < 1 {
		return []Todo{}
	}
	result := []Todo{}

	for _, todo := range tl.Todos {
		if !todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}

func (tl *TodoList) CountPending() int {
	if tl.Count() < 1 {
		return 0
	}

	count := 0
	for _, todo := range tl.Todos {
		if !todo.Completed {
			count++
		}
	}
	return count
}

func (tl *TodoList) GetCompleted() []Todo {
	result := []Todo{}

	for _, todo := range tl.Todos {
		if todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}

// Toggle completion status
func (tl *TodoList) Toggle(id int) error {
	i := tl.findIndex(id)
	if i < 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	tl.Todos[i].Completed = !tl.Todos[i].Completed
	return nil
}
