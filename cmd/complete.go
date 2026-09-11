package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var completeCmd = &cobra.Command{
	Use:   "complete [ID]",
	Short: "Complete a todo",
	Long:  "Complete a todo in the the todo list and mark status as completed",
	Args:  cobra.MinimumNArgs(1),
	RunE:  completeTodo,
}

func completeTodo(cmd *cobra.Command, args []string) error {
	todoId, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}
	if err := todoList.Complete(todoId); err != nil {
		return err
	}
	if err := saveTodos(); err != nil {
		return err
	}
	fmt.Printf("Completed todo with ID: %d\n", todoId)
	return nil
}
