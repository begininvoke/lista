package config

import (
	"path/filepath"
	"runtime"
	"testing"
)

func unixOnly(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		t.Skip("XDG_CONFIG_HOME is only honored on Linux/BSD/Unix")
	}
}

func TestBaseConfigDirWithXDG(t *testing.T) {
	unixOnly(t)

	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	got, err := baseConfigDir()
	if err != nil {
		t.Fatalf("baseConfigDir() error: %v", err)
	}

	want := filepath.Join(xdg, "lista")
	if got != want {
		t.Errorf("baseConfigDir() = %q, want %q", got, want)
	}
}

func TestBaseConfigDirFallback(t *testing.T) {
	unixOnly(t)

	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	got, err := baseConfigDir()
	if err != nil {
		t.Fatalf("baseConfigDir() error: %v", err)
	}

	want := filepath.Join(home, ".config", "lista")
	if got != want {
		t.Errorf("baseConfigDir() = %q, want %q", got, want)
	}
}

func TestFilePaths(t *testing.T) {
	unixOnly(t)

	tests := []struct {
		name string
		fn   func() (string, error)
		file string
	}{
		{"data", DataFilePath, "lista.json"},
		{"config", ConfigFilePath, "lista.config.json"},
	}

	t.Run("xdg", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				xdg := t.TempDir()
				t.Setenv("XDG_CONFIG_HOME", xdg)

				got, err := tt.fn()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				want := filepath.Join(xdg, "lista", tt.file)
				if got != want {
					t.Errorf("got %q, want %q", got, want)
				}
			})
		}
	})

	t.Run("home", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("XDG_CONFIG_HOME", "")
				t.Setenv("HOME", home)

				got, err := tt.fn()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				want := filepath.Join(home, ".config", "lista", tt.file)
				if got != want {
					t.Errorf("got %q, want %q", got, want)
				}
			})
		}
	})
}
