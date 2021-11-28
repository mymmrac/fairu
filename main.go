package main

import (
	"errors"
	"fmt"
	"github.com/muesli/reflow/wrap"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type level struct {
	path    string
	entries []os.DirEntry
}

func (l level) nextPath(i int) string {
	return filepath.Join(l.path, l.entries[i].Name())
}

type model struct {
	prevLvl *level
	currLvl *level
	nextLvl *level

	selected     int
	prevSelected int

	filter filterer

	width, height int

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

		for i, p := range m.prevLvl.entries {
			if p.Name() == filepath.Base(path) {
				m.prevSelected = i
				break
			}
		}
	} else {
		m.prevLvl = nil
	}

	if len(m.currLvl.entries) > 0 && m.currLvl.entries[m.selected].IsDir() {
		nextPath := m.currLvl.nextPath(m.selected)
		m.nextLvl, err = m.readDir(nextPath)
		if err != nil {
			m.err = err
			return m
		}
	} else {
		m.nextLvl = nil
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
		case key.Matches(msg, globalKeys.Quit) || key.Matches(msg, forceQuitKey):
			return m, tea.Quit
		case key.Matches(msg, globalKeys.NextEntry):
			m.selected++
			if m.selected >= len(m.currLvl.entries) {
				m.selected = 0
			}
			m = m.readLevels(m.currLvl.path)
		case key.Matches(msg, globalKeys.PrevEntry):
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.currLvl.entries) - 1
			}
			m = m.readLevels(m.currLvl.path)
		case key.Matches(msg, globalKeys.NextLvl):
			if m.currLvl.entries[m.selected].IsDir() {
				nextPath := m.currLvl.nextPath(m.selected)
				m.selected = 0
				m = m.readLevels(nextPath)
			}
		case key.Matches(msg, globalKeys.PrevLvl):
			if m.prevLvl != nil {
				prevPath := m.prevLvl.path
				m.selected = m.prevSelected
				m = m.readLevels(prevPath)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

var (
	borderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true)
	reversed    = lipgloss.NewStyle().Reverse(true)
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error occurred: %v\n\nPress any key to exit.\n", m.err)
	}

	s := strings.Builder{}

	prevS := strings.Builder{}
	currS := strings.Builder{}
	nextS := strings.Builder{}

	w := m.width - 2
	h := m.height - 2

	sidesWidth := int(float64(w) / 3.0)
	currOff := w - sidesWidth*3

	centerWidth := sidesWidth + currOff - 2

	if m.prevLvl != nil {
		prevS.WriteString(fmt.Sprintf("%s\n\n", wrap.String(m.prevLvl.path, sidesWidth)))

		for _, entry := range m.prevLvl.entries {
			prevS.WriteString(fmt.Sprintf("%s\n", wrap.String(entry.Name(), sidesWidth)))
		}
	}

	currS.WriteString(fmt.Sprintf("%s\n\n", wrap.String(m.currLvl.path, centerWidth)))

	for i, entry := range m.currLvl.entries {
		currFile := wrap.String(entry.Name(), centerWidth)
		if i == m.selected {
			currFile = reversed.Width(centerWidth).Render(currFile)
		}
		currS.WriteString(currFile + "\n")
	}

	if m.nextLvl != nil {
		nextS.WriteString(fmt.Sprintf("%s\n\n", wrap.String(m.nextLvl.path, sidesWidth)))

		for _, entry := range m.nextLvl.entries {
			nextS.WriteString(fmt.Sprintf("%s\n", wrap.String(entry.Name(), sidesWidth)))
		}
	}

	sidesStyle := lipgloss.NewStyle().
		Width(sidesWidth).Height(h).
		MaxHeight(h)

	centerStyle := lipgloss.NewStyle().
		Width(centerWidth).Height(h).
		MaxHeight(h).
		Border(lipgloss.NormalBorder(), false, true)

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left,
		sidesStyle.Render(prevS.String()),
		centerStyle.Render(currS.String()),
		sidesStyle.Render(nextS.String())))

	return borderStyle.Render(lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, s.String()))
}

func main() {
	if err := tea.NewProgram(newModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("Opps, something went wronk: %v", err)
		os.Exit(1)
	}
}
