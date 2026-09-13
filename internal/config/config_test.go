package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}
}

func TestDefaultTheme(t *testing.T) {
	want := Theme{
		Background:      "#3c3836",
		BackgroundAlt:   "#282828",
		TextPrimary:     "#ebdbb2",
		TextSecondary:   "#C5C7BC",
		TextMuted:       "#a89984",
		Error:           "#cc241d",
		Warning:         "#fe8019",
		Success:         "#b8bb26",
		Info:            "#83a598",
		PriorityHigh:    "#fb4934",
		PriorityMedium:  "#fe8019",
		PriorityLow:     "#b8bb26",
		Accent:          "#fabd2f",
		AccentSecondary: "#83a598",
		Border:          "#a89984",
	}

	got := DefaultTheme()
	if got != want {
		t.Errorf("DefaultTheme() = %+v, want %+v", got, want)
	}
}

func TestLoadConfigCreatesDefaultsToXDG(t *testing.T) {
	unixOnly(t)

	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	want := DefaultTheme()
	if cfg.Theme != want {
		t.Errorf("cfg.Theme = %+v, want %+v", cfg.Theme, want)
	}

	path := filepath.Join(xdg, "lista", "lista.config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("default config file not written: %v", err)
	}

	var written Config
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("written config is invalid JSON: %v", err)
	}
	if written.Theme != want {
		t.Errorf("written theme = %+v, want %+v", written.Theme, want)
	}
}

func TestLoadConfigCreatesDefaultsToHome(t *testing.T) {
	unixOnly(t)

	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	want := DefaultTheme()
	if cfg.Theme != want {
		t.Errorf("cfg.Theme = %+v, want %+v", cfg.Theme, want)
	}

	path := filepath.Join(home, ".config", "lista", "lista.config.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("default config file not written to HOME: %v", err)
	}
}

func TestLoadConfigExistingPartial(t *testing.T) {
	unixOnly(t)

	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	custom := Theme{Background: "#000000"}
	data, err := json.Marshal(Config{Theme: custom})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	writeConfigFile(t, filepath.Join(xdg, "lista", "lista.config.json"), data)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	if cfg.Theme.Background != "#000000" {
		t.Errorf("Background = %q, want %q", cfg.Theme.Background, "#000000")
	}

	def := DefaultTheme()
	if cfg.Theme.TextPrimary != def.TextPrimary {
		t.Errorf("TextPrimary = %q, want default %q", cfg.Theme.TextPrimary, def.TextPrimary)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	unixOnly(t)

	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	writeConfigFile(t, filepath.Join(xdg, "lista", "lista.config.json"), []byte("not json"))

	if _, err := LoadConfig(); err == nil {
		t.Fatal("LoadConfig() expected error for invalid JSON, got nil")
	}
}

func TestMergeWithDefaults(t *testing.T) {
	got := mergeWithDefaults(Theme{})
	want := DefaultTheme()
	if got != want {
		t.Errorf("mergeWithDefaults(Theme{}) = %+v, want %+v", got, want)
	}
}

func TestMergeWithDefaultsKeepsCustom(t *testing.T) {
	got := mergeWithDefaults(Theme{
		Background: "#ffffff",
		Accent:     "#00ff00",
		Border:     "#000000",
	})
	want := DefaultTheme()
	want.Background = "#ffffff"
	want.Accent = "#00ff00"
	want.Border = "#000000"

	if got != want {
		t.Errorf("mergeWithDefaults() = %+v, want %+v", got, want)
	}
}
