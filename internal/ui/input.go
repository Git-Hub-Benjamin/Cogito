package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	textInput textinput.Model
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
	return m, cmd
}

func (m InputModel) View() string {
	// Re-apply accent color each render so it stays in sync
	m.textInput.PromptStyle = InputPromptStyle
	return m.textInput.View()
}

func (m InputModel) Value() string {
	return m.textInput.Value()
}

func (m *InputModel) SetValue(s string) {
	m.textInput.SetValue(s)
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
