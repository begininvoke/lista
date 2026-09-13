package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kwame-Owusu/lista/internal/storage"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports todos",
	Long:  "Exports todos to stdout in .md syntax",
	Args:  cobra.MinimumNArgs(0),
	RunE:  exportTodos,
}

func exportTodos(cmd *cobra.Command, args []string) error {
	todos, err := storage.LoadTodos(dataFile)
	if err != nil {
		return fmt.Errorf("reading todos: %w", err)
	}

	var sb strings.Builder

	sb.WriteString("# Lista\n\n")
	sb.WriteString("## Todos\n")
	for _, todo := range todos {
		checked := " "
		if todo.Completed {
			checked = "x"
		}
		details := fmt.Sprintf("(%s", todo.Priority)
		if !todo.CreatedAt.IsZero() {
			details += fmt.Sprintf(" %s", todo.CreatedAt.Format("2006-01-02"))
		}
		details += ")"
		fmt.Fprintf(&sb, "- [%s] %s %s\n", checked, todo.Title, details)
	}

	hasNotes := false
	for _, todo := range todos {
		if strings.TrimSpace(todo.Notes) != "" {
			hasNotes = true
			break
		}
	}

	if hasNotes {
		sb.WriteString("\n## Notes\n")
		for _, todo := range todos {
			if strings.TrimSpace(todo.Notes) == "" {
				continue
			}
			fmt.Fprintf(&sb, "- **%s**\n  %s\n", todo.Title, todo.Notes)
		}
	}

	fmt.Fprint(os.Stdout, sb.String())
	return nil
}
