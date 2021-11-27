package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type level struct {
	path    string
	entries []os.DirEntry
}

type model struct {
	prevLvl *level
	currLvl *level
	nextLvl *level

	currentSelected int

	filter filterer

	err error
}

const maxArgs = 2

func (m model) readDir(path string) (*level, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	lvl := &level{
		path:    path,
		entries: []os.DirEntry{},
	}

	for _, entry := range entries {
		if m.filter.filter(path, entry) {
			lvl.entries = append(lvl.entries, entry)
		}
	}

	return lvl, nil
}

func (m model) readLevels(path string) model {
	var err error

	m.currLvl, err = m.readDir(path)
	if err != nil {
		return model{err: err}
	}

	prevPath := filepath.Dir(path)
	if prevPath != path {
		m.prevLvl, err = m.readDir(prevPath)
		if err != nil {
			m.err = err
			return m
		}
	}

	var nextPath string
	for i, entry := range m.currLvl.entries {
		if entry.IsDir() {
			nextPath = filepath.Join(path, entry.Name())
			m.currentSelected = i
			break
		}
	}

	if nextPath != "" {
		m.nextLvl, err = m.readDir(nextPath)
		if err != nil {
			m.err = err
			return m
		}
	}

	return m
}

func newModel() *model {
	m := model{
		filter: dotFilter(false),
	}

	if len(os.Args) > maxArgs {
		m.err = errors.New("too many arguments")
		return &m
	}

	localPath := "."
	if len(os.Args) == maxArgs {
		localPath = os.Args[1]
	}

	path, err := filepath.Abs(localPath)
	if err != nil {
		m.err = err
		return &m
	}

	info, err := os.Stat(path)
	if err != nil {
		m.err = err
		return &m
	}

	if !info.IsDir() {
		m.err = fmt.Errorf("%s is a file, expected directory", path)
		return &m
	}

	m = m.readLevels(path)

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

	if m.prevLvl != nil {
		s.WriteString(fmt.Sprintf("Previous path: %s\n", m.prevLvl.path))

		for _, entry := range m.prevLvl.entries {
			s.WriteString(fmt.Sprintf("%s\n", entry.Name()))
		}
	}

	s.WriteString(fmt.Sprintf("Current path: %s\n", m.currLvl.path))

	s.WriteString("Current files:\n")
	for _, entry := range m.currLvl.entries {
		s.WriteString(fmt.Sprintf("%s\n", entry.Name()))
	}

	if m.nextLvl != nil {
		s.WriteString(fmt.Sprintf("Next path: %s\n", m.nextLvl.path))

		for _, entry := range m.nextLvl.entries {
			s.WriteString(fmt.Sprintf("%s\n", entry.Name()))
		}
	}

	return s.String()
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Opps, something went wronk: %v", err)
		os.Exit(1)
	}
}
