package main

import "github.com/charmbracelet/bubbles/key"

var forceQuitKey = key.NewBinding(
	key.WithKeys("ctrl+c"),
)

type globalKeybindings struct {
	Quit,
	PrevLvl, NextLvl,
	PrevEntry, NextEntry,
	_ key.Binding
}

var globalKeys = globalKeybindings{
	Quit: key.NewBinding(
		key.WithKeys("q"),
	),

	PrevLvl: key.NewBinding(
		key.WithKeys("left"),
	),

	NextLvl: key.NewBinding(
		key.WithKeys("right"),
	),

	PrevEntry: key.NewBinding(
		key.WithKeys("up"),
	),

	NextEntry: key.NewBinding(
		key.WithKeys("down"),
	),
}
