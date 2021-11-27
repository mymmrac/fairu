package main

import "github.com/charmbracelet/bubbles/key"

type globalKeybindings struct {
	Quit key.Binding
}

var globalKeys = globalKeybindings{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
	),
}
