package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a todo",
	Long:  "Delete a todo and remove it from the list",
	Args:  cobra.MinimumNArgs(1),
	RunE:  deleteTodo,
}

func deleteTodo(cmd *cobra.Command, args []string) error {
	todoId, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}
	if err := todoList.Delete(todoId); err != nil {
		return err
	}
	if err := saveTodos(); err != nil {
		return err
	}
	fmt.Printf("Deleted todo with ID: %d\n", todoId)
	return nil
}
