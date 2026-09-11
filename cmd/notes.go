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

	newNotes := strings.Join(args[1:], " ")
	if err := todoList.AppendNotes(id, newNotes); err != nil {
		return err
	}
	return saveTodos()
}
