package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "Edit the title of a todo",
	Long:  "Edit the title of a todo, given the correct ID",
	Args:  cobra.MinimumNArgs(2),
	RunE:  editTodoTitle,
}

func editTodoTitle(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}

	editTitle := strings.Join(args[1:], " ")
	if err := todoList.Edit(id, editTitle); err != nil {
		return fmt.Errorf("editing todo %d: %w", id, err)
	}

	return saveTodos()
}
