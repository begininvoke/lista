package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var addNotesCmd = &cobra.Command{
	Use:   "notes [id]",
	Short: "Adds notes to a todo",
	Long:  "Adds notes to a todo, if not present. if present it appends to the existing note",
	Args:  cobra.MinimumNArgs(1),
	RunE:  addNotes,
}

func addNotes(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid todo ID %q", args[0])
	}

	todo, err := todoList.GetByID(id)
	if err != nil {
		return err
	}
	newNotes := strings.Join(args[1:], " ")

	if todo.Notes != "" {
		todo.Notes += " " + newNotes
	} else {
		todo.Notes = newNotes
	}
	return saveTodos()
}
