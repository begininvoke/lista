package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/config"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/storage"
	"github.com/kwame-Owusu/lista/internal/tui"
	"github.com/spf13/cobra"
)

var version = "dev"
var todoList *models.TodoList
var dataFile string //$HOME/.config/lista, where our json configs live

var styleOnce sync.Once

func loadStyles() {
	styleOnce.Do(func() {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v, using defaults\n", err)
			cfg = &config.Config{Theme: config.DefaultTheme()}
		}
		tui.InitStyles(cfg.Theme)
	})
}

func loadTodos() error {
	path, err := config.DataFilePath()
	if err != nil {
		return fmt.Errorf("resolving config path: %w", err)
	}
	dataFile = path

	// 0755 = rwx for owner, rx for group and others.
	// This is the standard permission set for config directories:
	// - Owner can read/write config files
	// - Others can traverse the directory but not modify its contents
	permissions := 0755

	if err := os.MkdirAll(filepath.Dir(dataFile), os.FileMode(permissions)); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	todos, err := storage.LoadTodos(dataFile)
	if err != nil {
		// First run: no data file yet, start fresh.
		if errors.Is(err, fs.ErrNotExist) {
			todoList = models.NewTodoList()
			return nil
		}

		// The file exists but couldn't be loaded (corrupt JSON, permissions,
		// ...). Move it aside so the next save doesn't silently destroy it.
		fmt.Fprintf(os.Stderr, "Warning: unable to read %s (%v)\n", dataFile, err)
		backup, backupErr := storage.BackupCorruptFile(dataFile)
		if backupErr != nil {
			fmt.Fprintf(os.Stderr, "Error: could not back up the file: %v\n", backupErr)
		} else {
			fmt.Fprintf(os.Stderr, "Your data file was moved to %s; starting with an empty list.\n", backup)
		}

		todoList = models.NewTodoList()
		return nil
	}

	// File exists - create TodoList and populate it
	todoList = models.NewTodoList()
	for _, todo := range todos {
		todoList.Todos = append(todoList.Todos, todo)
		// Update NextID to be higher than highest existing ID
		if todo.ID >= todoList.NextID {
			todoList.NextID = todo.ID + 1
		}
	}
	return nil
}

func saveTodos() error {
	err := storage.SaveTodos(todoList.Todos, dataFile)
	if err != nil {
		return fmt.Errorf("saving todos: %w", err)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:           "lista",
	Short:         "A minimal todo CLI program",
	Long:          `Lista is a simple and aesthetic CLI app to manage your todos on the terminal.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadTodos()
	},
	Run: func(cmd *cobra.Command, args []string) {
		loadStyles()
		m := tui.NewModel(todoList, dataFile)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(uncompleteCmd)
	rootCmd.AddCommand(toggleCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(addNotesCmd)
	rootCmd.AddCommand(exportCmd)
}
