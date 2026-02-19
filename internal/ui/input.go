package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var suggestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))

type InputModel struct {
	textInput  textinput.Model
	suggestion string // ghost text shown after cursor
}

func NewInputModel() InputModel {
	ti := textinput.New()
	ti.Placeholder = "Ask anything..."
	ti.PromptStyle = InputPromptStyle
	ti.Prompt = "❯ "
	ti.Focus()
	ti.CharLimit = 2000
	return InputModel{textInput: ti}
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	// Clear suggestion on any keypress (tab handling is in app.go)
	m.suggestion = ""
	return m, cmd
}

func (m InputModel) View() string {
	m.textInput.PromptStyle = InputPromptStyle
	view := m.textInput.View()
	if m.suggestion != "" {
		view += suggestionStyle.Render(m.suggestion)
	}
	return view
}

func (m InputModel) Value() string {
	return m.textInput.Value()
}

func (m *InputModel) SetValue(s string) {
	m.textInput.SetValue(s)
	// Move cursor to end of the new value
	m.textInput.SetCursor(len(s))
}

func (m *InputModel) SetSuggestion(s string) {
	m.suggestion = s
}

func (m *InputModel) Focus() tea.Cmd {
	return m.textInput.Focus()
}

func (m *InputModel) Blur() {
	m.textInput.Blur()
}

func (m *InputModel) SetWidth(w int) {
	m.textInput.Width = w
}
