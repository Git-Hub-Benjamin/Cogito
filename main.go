package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/benji/cogito/internal/app"
	"github.com/benji/cogito/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	m := app.NewModel(cfg)

	// When rendering at top without clearing, move cursor to top-left
	// so Bubble Tea's inline renderer starts from position (1,1).
	// Old terminal content below the box stays visible.
	if !cfg.ClearScreen && cfg.Position == "top" {
		fmt.Print("\033[H")
	}

	opts := []tea.ProgramOption{}
	if cfg.ClearScreen {
		opts = append(opts, tea.WithAltScreen())
	}
	p := tea.NewProgram(m, opts...)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
