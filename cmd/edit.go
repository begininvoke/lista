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
	Run:   editTodoTitle,
}

func editTodoTitle(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error occurred converting id to int: %s\n", err)
		return
	}

	editTitle := strings.Join(args[1:], " ")
	err = todoList.Edit(id, editTitle)
	if err != nil {
		fmt.Printf("Error editing todo with id: %d, and string: %s, %s\n", id, editTitle, err)
		return
	}

	saveTodos()
}
