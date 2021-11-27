package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	path      string
	files     []os.FileInfo
	fileNames []string
	filter    filterer

	err error
}

const maxArgs = 2

func (m model) readDir() model {
	m.files, m.err = ioutil.ReadDir(m.path)
	if m.err != nil {
		return m
	}

	m.fileNames = []string{}
	for _, fi := range m.files {
		if m.filter.filter(m.path, fi) {
			m.fileNames = append(m.fileNames, fi.Name())
		}
	}
	return m
}

func newModel() *model {
	if len(os.Args) > maxArgs {
		return &model{err: errors.New("too many arguments")}
	}

	path := "."
	if len(os.Args) == maxArgs {
		path = os.Args[1]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return &model{err: err}
	}

	info, err := os.Stat(path)
	if err != nil {
		return &model{err: err}
	}

	if !info.IsDir() {
		return &model{err: fmt.Errorf("%s is a file, expected directory", absPath)}
	}

	m := model{
		path:   absPath,
		filter: dotFilter(false),
	}

	m = m.readDir()

	return &m
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

	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("Current path: %s\n", m.path))

	s.WriteString(fmt.Sprintf("Files:\n%s\n", strings.Join(m.fileNames, "\n")))

	return s.String()
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Opps, something went wronk: %v", err)
		os.Exit(1)
	}
}
