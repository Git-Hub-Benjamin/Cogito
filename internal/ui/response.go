package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResponseModel struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func NewResponseModel() ResponseModel {
	return ResponseModel{}
}

func (m *ResponseModel) SetSize(width, height int) {
	if !m.ready {
		m.viewport = viewport.New(width, height)
		m.viewport.Style = lipgloss.NewStyle()
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = height
	}
	m.viewport.SetContent(m.content)
}

func (m ResponseModel) Update(msg tea.Msg) (ResponseModel, tea.Cmd) {
	if !m.ready {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ResponseModel) View() string {
	if !m.ready {
		return ""
	}
	return m.viewport.View()
}

func (m *ResponseModel) AppendContent(chunk string) {
	m.content += chunk
	if m.ready {
		m.viewport.SetContent(m.content)
		m.viewport.GotoBottom()
	}
}

func (m *ResponseModel) Clear() {
	m.content = ""
	if m.ready {
		m.viewport.SetContent("")
		m.viewport.GotoTop()
	}
}

func (m ResponseModel) Content() string {
	return m.content
}
