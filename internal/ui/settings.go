package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsSavedMsg struct {
	APIKey       string
	BaseURL      string
	DefaultModel string
	ClearScreen  bool
	Position     string
	AccentColor  string
}

// settingsItem is either a group header or a text input field.
type settingsItem struct {
	label    string
	isGroup  bool   // true = collapsible header, not an input
	groupID  string // which group this input belongs to ("api" or "")
	inputIdx int    // index into SettingsModel.inputs, -1 for group headers
}

type SettingsModel struct {
	inputs   []textinput.Model
	items    []settingsItem
	cursor   int
	width    int
	expanded map[string]bool // which groups are expanded
}

const (
	inputAPIKey = iota
	inputBaseURL
	inputModel
	inputClearScreen
	inputPosition
	inputAccentColor
	inputCount
)

func NewSettingsModel(apiKey, baseURL, defaultModel string, clearScreen bool, position, accentColor string) SettingsModel {
	inputs := make([]textinput.Model, inputCount)

	inputs[inputAPIKey] = textinput.New()
	inputs[inputAPIKey].Placeholder = "sk-..."
	inputs[inputAPIKey].EchoMode = textinput.EchoPassword
	inputs[inputAPIKey].SetValue(apiKey)
	inputs[inputAPIKey].CharLimit = 256
	inputs[inputAPIKey].Width = 50

	inputs[inputBaseURL] = textinput.New()
	inputs[inputBaseURL].Placeholder = "leave empty for default OpenAI"
	inputs[inputBaseURL].SetValue(baseURL)
	inputs[inputBaseURL].CharLimit = 256
	inputs[inputBaseURL].Width = 50

	inputs[inputModel] = textinput.New()
	inputs[inputModel].Placeholder = "gpt-4o-mini"
	inputs[inputModel].SetValue(defaultModel)
	inputs[inputModel].CharLimit = 100
	inputs[inputModel].Width = 50

	clearVal := "no"
	if clearScreen {
		clearVal = "yes"
	}
	inputs[inputClearScreen] = textinput.New()
	inputs[inputClearScreen].Placeholder = "yes/no"
	inputs[inputClearScreen].SetValue(clearVal)
	inputs[inputClearScreen].CharLimit = 3
	inputs[inputClearScreen].Width = 50

	if position == "" {
		position = "bottom"
	}
	inputs[inputPosition] = textinput.New()
	inputs[inputPosition].Placeholder = "top/bottom"
	inputs[inputPosition].SetValue(position)
	inputs[inputPosition].CharLimit = 6
	inputs[inputPosition].Width = 50

	if accentColor == "" {
		accentColor = "#FF6F61"
	}
	inputs[inputAccentColor] = textinput.New()
	inputs[inputAccentColor].Placeholder = "#FF6F61"
	inputs[inputAccentColor].SetValue(accentColor)
	inputs[inputAccentColor].CharLimit = 7
	inputs[inputAccentColor].Width = 50

	items := []settingsItem{
		{label: "API Configuration", isGroup: true, groupID: "api", inputIdx: -1},
		{label: "API Key", groupID: "api", inputIdx: inputAPIKey},
		{label: "Base URL (Groq, OpenRouter, Ollama...)", groupID: "api", inputIdx: inputBaseURL},
		{label: "Default Model", groupID: "api", inputIdx: inputModel},
		{label: "Clear Screen (yes/no)", inputIdx: inputClearScreen},
		{label: "Position (top/bottom)", inputIdx: inputPosition},
		{label: "Accent Color (hex)", inputIdx: inputAccentColor},
	}

	return SettingsModel{
		inputs:   inputs,
		items:    items,
		cursor:   0,
		expanded: map[string]bool{"api": false},
	}
}

func (m SettingsModel) visibleItems() []int {
	var visible []int
	for i, item := range m.items {
		if item.groupID != "" && !item.isGroup {
			if !m.expanded[item.groupID] {
				continue
			}
		}
		visible = append(visible, i)
	}
	return visible
}

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	visible := m.visibleItems()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			item := m.items[visible[m.cursor]]
			if !item.isGroup && item.inputIdx >= 0 {
				m.inputs[item.inputIdx].Blur()
			}
			m.cursor = (m.cursor + 1) % len(visible)
			return m, m.focusCurrent()

		case "shift+tab", "up":
			item := m.items[visible[m.cursor]]
			if !item.isGroup && item.inputIdx >= 0 {
				m.inputs[item.inputIdx].Blur()
			}
			m.cursor = (m.cursor - 1 + len(visible)) % len(visible)
			return m, m.focusCurrent()

		case "enter":
			item := m.items[visible[m.cursor]]
			if item.isGroup {
				m.expanded[item.groupID] = !m.expanded[item.groupID]
				return m, nil
			}
			// Last visible item → save
			if m.cursor == len(visible)-1 {
				return m, m.buildSaveMsg()
			}
			// Otherwise advance to next
			if item.inputIdx >= 0 {
				m.inputs[item.inputIdx].Blur()
			}
			m.cursor++
			if m.cursor >= len(visible) {
				m.cursor = 0
			}
			return m, m.focusCurrent()
		}
	}

	// Forward to active input
	item := m.items[visible[m.cursor]]
	if !item.isGroup && item.inputIdx >= 0 {
		var cmd tea.Cmd
		m.inputs[item.inputIdx], cmd = m.inputs[item.inputIdx].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *SettingsModel) focusCurrent() tea.Cmd {
	visible := m.visibleItems()
	item := m.items[visible[m.cursor]]
	if !item.isGroup && item.inputIdx >= 0 {
		return m.inputs[item.inputIdx].Focus()
	}
	return nil
}

func (m SettingsModel) buildSaveMsg() tea.Cmd {
	clearVal := strings.TrimSpace(strings.ToLower(m.inputs[inputClearScreen].Value()))
	clearScreen := clearVal == "yes" || clearVal == "y" || clearVal == "true"

	pos := strings.TrimSpace(strings.ToLower(m.inputs[inputPosition].Value()))
	if pos != "top" {
		pos = "bottom"
	}
	color := strings.TrimSpace(m.inputs[inputAccentColor].Value())

	return func() tea.Msg {
		return SettingsSavedMsg{
			APIKey:       m.inputs[inputAPIKey].Value(),
			BaseURL:      strings.TrimSpace(m.inputs[inputBaseURL].Value()),
			DefaultModel: m.inputs[inputModel].Value(),
			ClearScreen:  clearScreen,
			Position:     pos,
			AccentColor:  color,
		}
	}
}

func (m SettingsModel) View() string {
	visible := m.visibleItems()

	var b strings.Builder
	b.WriteString(TitleStyle.Render("Settings") + "\n\n")

	for vi, itemIdx := range visible {
		item := m.items[itemIdx]
		isCursor := vi == m.cursor

		if item.isGroup {
			arrow := "▸"
			if m.expanded[item.groupID] {
				arrow = "▾"
			}
			if isCursor {
				b.WriteString(SelectedStyle.Render(fmt.Sprintf("%s %s", arrow, item.label)) + "\n")
			} else {
				b.WriteString(DimStyle.Render(fmt.Sprintf("%s %s", arrow, item.label)) + "\n")
			}
			if !m.expanded[item.groupID] {
				// Show a summary of current values
				model := m.inputs[inputModel].Value()
				url := m.inputs[inputBaseURL].Value()
				if url == "" {
					url = "openai"
				}
				b.WriteString(DimStyle.Render(fmt.Sprintf("    %s @ %s", model, url)) + "\n")
			}
			b.WriteString("\n")
			continue
		}

		indent := "  "
		if item.groupID != "" {
			indent = "    "
		}

		if isCursor {
			b.WriteString(SelectedStyle.Render(fmt.Sprintf("%s▸ %s:", indent, item.label)) + "\n")
		} else {
			b.WriteString(DimStyle.Render(fmt.Sprintf("%s  %s:", indent, item.label)) + "\n")
		}
		b.WriteString(indent + "  " + m.inputs[item.inputIdx].View() + "\n\n")
	}

	b.WriteString(DimStyle.Render("↑↓ navigate • enter expand/save • esc back"))
	return b.String()
}

func (m *SettingsModel) SetWidth(w int) {
	m.width = w
	for i := range m.inputs {
		m.inputs[i].Width = w - 10
	}
}
