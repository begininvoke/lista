package cmd

import (
	"fmt"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/tui"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

var (
	listPendingFlag  bool
	listDoneFlag     bool
	listPriorityFlag string
	listSearchFlag   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos",
	Long:  "List all todo tasks in the todo list",
	Args:  cobra.MinimumNArgs(0),
	RunE:  listTodos,
}

func init() {
	listCmd.Flags().BoolVar(&listPendingFlag, "pending", false, "Show only pending todos")
	listCmd.Flags().BoolVarP(&listDoneFlag, "done", "d", false, "Show only completed todos")
	listCmd.Flags().StringVarP(&listPriorityFlag, "priority", "p", "", "Filter by priority (high/h, medium/m, low/l)")
	listCmd.Flags().StringVarP(&listSearchFlag, "search", "s", "", "Filter by title or notes (case-insensitive)")
}

type listFilter struct {
	pending  bool
	done     bool
	priority *models.Priority
	search   string
}

func listTodos(cmd *cobra.Command, args []string) error {
	if listPendingFlag && listDoneFlag {
		return fmt.Errorf("--pending and --done cannot be used together")
	}

	filter := listFilter{
		pending: listPendingFlag,
		done:    listDoneFlag,
		search:  strings.TrimSpace(listSearchFlag),
	}

	if listPriorityFlag != "" {
		priority, err := models.ParsePriority(listPriorityFlag)
		if err != nil {
			return fmt.Errorf("filtering by priority: %w", err)
		}
		filter.priority = &priority
	}

	loadStyles()
	todos := filterTodos(todoList.List(), filter)

	sort.Slice(todos, func(i, j int) bool {
		if todos[i].Completed != todos[j].Completed {
			return !todos[i].Completed
		}
		return todos[i].Priority > todos[j].Priority
	})

	// Header
	fmt.Printf("%-4s %-10s %-10s %-14s %s\n",
		tui.RenderHeader("ID"),
		tui.RenderHeader("STATUS"),
		tui.RenderHeader("PRIORITY"),
		tui.RenderHeader("TIME"),
		tui.RenderHeader("TITLE"),
	)
	fmt.Println(tui.RenderMuted(strings.Repeat("-", 80)))

	for _, todo := range todos {
		title := todo.Title
		if len(todo.Notes) > 0 {
			title += " *"
		}

		fmt.Println(
			fmt.Sprint(todo.ID),
			tui.RenderStatus(todo.Completed),
			tui.RenderPriority(todo.Priority.String()),
			tui.RenderTimestamp(todo.TimeAgo()),
			tui.RenderTodoTitle(title, todo.Completed),
		)
	}

	return nil
}

func filterTodos(todos []models.Todo, f listFilter) []models.Todo {
	result := make([]models.Todo, 0, len(todos))

	for _, todo := range todos {
		if f.pending && todo.Completed {
			continue
		}
		if f.done && !todo.Completed {
			continue
		}
		if f.priority != nil && todo.Priority != *f.priority {
			continue
		}
		if f.search != "" &&
			!strings.Contains(strings.ToLower(todo.Title), strings.ToLower(f.search)) &&
			!strings.Contains(strings.ToLower(todo.Notes), strings.ToLower(f.search)) {
			continue
		}
		result = append(result, todo)
	}

	return result
}