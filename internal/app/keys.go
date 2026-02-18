package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Submit      key.Binding
	Quit        key.Binding
	FocusInput  key.Binding
	TabComplete key.Binding
	Cancel      key.Binding
}

var keys = keyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
	FocusInput: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "focus input"),
	),
	TabComplete: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "autocomplete"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
