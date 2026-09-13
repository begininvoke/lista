package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle [ID]",
	Short: "Toggle a todo's completion status",
	Long:  "Toggle a todo's completion status between pending and completed",
	Args:  cobra.MinimumNArgs(1),
	RunE:  toggleTodo,
}

func toggleTodo(cmd *cobra.Command, args []string) error {
	todoId, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}
	if err := todoList.Toggle(todoId); err != nil {
		return err
	}
	if err := saveTodos(); err != nil {
		return err
	}

	todo, err := todoList.GetByID(todoId)
	if err != nil {
		return err
	}
	status := "pending"
	if todo.Completed {
		status = "completed"
	}
	fmt.Printf("Toggled todo with ID: %d (%s)\n", todoId, status)
	return nil
}
