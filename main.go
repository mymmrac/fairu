package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	path string
	err  error
}

func newModel() *model {
	path := "."
	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return &model{err: err}
	}

	exist, err := exists(absPath)
	if err != nil {
		return &model{err: err}
	}

	if !exist {
		return &model{err: fmt.Errorf("file %s doesn't exist", absPath)}
	}

	return &model{
		path: absPath,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, tea.Quit
		}

		return m, tea.ExitAltScreen
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, globalKeys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error occurred: %v\n\nPress any key to exit.\n", m.err)
	}

	return fmt.Sprintf("Current path: %s\n", m.path)
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Opps, something went wronk: %v", err)
		os.Exit(1)
	}
}
