package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var uncompleteCmd = &cobra.Command{
	Use:   "uncomplete [ID]",
	Short: "Uncomplete a todo",
	Long:  "Uncomplete a todo in the todo list and mark status as pending",
	Args:  cobra.MinimumNArgs(1),
	RunE:  uncompleteTodo,
}

func uncompleteTodo(cmd *cobra.Command, args []string) error {
	todoId, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}
	if err := todoList.Uncomplete(todoId); err != nil {
		return err
	}
	if err := saveTodos(); err != nil {
		return err
	}
	fmt.Printf("Uncompleted todo with ID: %d\n", todoId)
	return nil
}
