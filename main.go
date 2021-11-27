package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func newModel() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return "Hello World!\n"
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Opps, something went wronk: %v", err)
		os.Exit(1)
	}
}
