package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsSavedMsg struct {
	APIKey             string
	BaseURL            string
	DefaultModel       string
	ClearScreen        bool
	Position           string
	AccentColor        string
	CustomInstructions string
	IncludeCWD         bool
	MaxResponseLines   int
}

// settingsItem is either a group header, a text input field, or a save button.
type settingsItem struct {
	label     string
	hint      string // optional hint shown below the label
	isGroup   bool   // true = collapsible header, not an input
	isSaveBtn bool   // true = save button
	groupID   string // which group this input belongs to ("api", "prompt", "display")
	inputIdx  int    // index into SettingsModel.inputs, -1 for group headers/buttons
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
	inputCustomInstructions
	inputIncludeCWD
	inputMaxResponseLines
	inputClearScreen
	inputPosition
	inputAccentColor
	inputCount
)

func NewSettingsModel(apiKey, baseURL, defaultModel, customInstructions string, includeCWD, clearScreen bool, position, accentColor string, maxResponseLines int) SettingsModel {
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

	inputs[inputCustomInstructions] = textinput.New()
	inputs[inputCustomInstructions].Placeholder = "e.g. Always respond in Python, be extra brief..."
	inputs[inputCustomInstructions].SetValue(customInstructions)
	inputs[inputCustomInstructions].CharLimit = 500
	inputs[inputCustomInstructions].Width = 50

	cwdVal := "no"
	if includeCWD {
		cwdVal = "yes"
	}
	inputs[inputIncludeCWD] = textinput.New()
	inputs[inputIncludeCWD].Placeholder = "yes/no"
	inputs[inputIncludeCWD].SetValue(cwdVal)
	inputs[inputIncludeCWD].CharLimit = 3
	inputs[inputIncludeCWD].Width = 50

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

	if maxResponseLines <= 0 {
		maxResponseLines = 8
	}
	inputs[inputMaxResponseLines] = textinput.New()
	inputs[inputMaxResponseLines].Placeholder = "8"
	inputs[inputMaxResponseLines].SetValue(fmt.Sprintf("%d", maxResponseLines))
	inputs[inputMaxResponseLines].CharLimit = 3
	inputs[inputMaxResponseLines].Width = 50

	items := []settingsItem{
		{label: "API Configuration", isGroup: true, groupID: "api", inputIdx: -1},
		{label: "API Key", groupID: "api", inputIdx: inputAPIKey},
		{label: "Base URL (Groq, OpenRouter, Ollama...)", groupID: "api", inputIdx: inputBaseURL},
		{label: "Default Model", groupID: "api", inputIdx: inputModel},

		{label: "Prompt & Context", isGroup: true, groupID: "prompt", inputIdx: -1},
		{label: "Custom Instructions", groupID: "prompt", inputIdx: inputCustomInstructions},
		{label: "Send Directory Context (yes/no)", hint: "Sends your current working directory to the model for relevant answers", groupID: "prompt", inputIdx: inputIncludeCWD},

		{label: "Display", isGroup: true, groupID: "display", inputIdx: -1},
		{label: "Max Response Lines", hint: "Lines shown before pager activates (default: 8)", groupID: "display", inputIdx: inputMaxResponseLines},
		{label: "Clear Screen (yes/no)", groupID: "display", inputIdx: inputClearScreen},
		{label: "Position (top/bottom)", groupID: "display", inputIdx: inputPosition},
		{label: "Accent Color (hex)", groupID: "display", inputIdx: inputAccentColor},

		{label: "Save & Exit", isSaveBtn: true, inputIdx: -1},
	}

	return SettingsModel{
		inputs:   inputs,
		items:    items,
		cursor:   0,
		expanded: map[string]bool{"api": false, "prompt": false, "display": false},
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
			if item.isSaveBtn {
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
	cwdVal := strings.TrimSpace(strings.ToLower(m.inputs[inputIncludeCWD].Value()))
	includeCWD := cwdVal == "yes" || cwdVal == "y" || cwdVal == "true"

	maxLines, err := strconv.Atoi(strings.TrimSpace(m.inputs[inputMaxResponseLines].Value()))
	if err != nil || maxLines < 3 {
		maxLines = 8
	}

	clearVal := strings.TrimSpace(strings.ToLower(m.inputs[inputClearScreen].Value()))
	clearScreen := clearVal == "yes" || clearVal == "y" || clearVal == "true"

	pos := strings.TrimSpace(strings.ToLower(m.inputs[inputPosition].Value()))
	if pos != "top" {
		pos = "bottom"
	}
	color := strings.TrimSpace(m.inputs[inputAccentColor].Value())

	return func() tea.Msg {
		return SettingsSavedMsg{
			APIKey:             m.inputs[inputAPIKey].Value(),
			BaseURL:            strings.TrimSpace(m.inputs[inputBaseURL].Value()),
			DefaultModel:       m.inputs[inputModel].Value(),
			CustomInstructions: strings.TrimSpace(m.inputs[inputCustomInstructions].Value()),
			IncludeCWD:         includeCWD,
			MaxResponseLines:   maxLines,
			ClearScreen:        clearScreen,
			Position:           pos,
			AccentColor:        color,
		}
	}
}

// groupSummary returns a short summary string for a collapsed group.
func (m SettingsModel) groupSummary(groupID string) string {
	switch groupID {
	case "api":
		model := m.inputs[inputModel].Value()
		url := m.inputs[inputBaseURL].Value()
		if url == "" {
			url = "openai"
		}
		return fmt.Sprintf("%s @ %s", model, url)
	case "prompt":
		parts := []string{}
		if ci := m.inputs[inputCustomInstructions].Value(); ci != "" {
			if len(ci) > 30 {
				ci = ci[:27] + "..."
			}
			parts = append(parts, fmt.Sprintf("\"%s\"", ci))
		} else {
			parts = append(parts, "no custom instructions")
		}
		cwdVal := strings.TrimSpace(strings.ToLower(m.inputs[inputIncludeCWD].Value()))
		if cwdVal == "yes" || cwdVal == "y" || cwdVal == "true" {
			parts = append(parts, "dir context on")
		} else {
			parts = append(parts, "dir context off")
		}
		return strings.Join(parts, " • ")
	case "display":
		clearVal := strings.TrimSpace(strings.ToLower(m.inputs[inputClearScreen].Value()))
		pos := m.inputs[inputPosition].Value()
		color := m.inputs[inputAccentColor].Value()
		lines := m.inputs[inputMaxResponseLines].Value()
		clear := "no"
		if clearVal == "yes" || clearVal == "y" || clearVal == "true" {
			clear = "yes"
		}
		return fmt.Sprintf("%s lines • clear: %s • %s • %s", lines, clear, pos, color)
	}
	return ""
}

func (m SettingsModel) View() string {
	visible := m.visibleItems()

	var b strings.Builder
	b.WriteString(TitleStyle.Render("Settings") + "\n\n")

	for vi, itemIdx := range visible {
		item := m.items[itemIdx]
		isCursor := vi == m.cursor

		if item.isSaveBtn {
			if isCursor {
				b.WriteString(SelectedStyle.Render("  ▸ [ "+item.label+" ]") + "\n")
			} else {
				b.WriteString(DimStyle.Render("    [ "+item.label+" ]") + "\n")
			}
			continue
		}

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
				summary := m.groupSummary(item.groupID)
				if summary != "" {
					b.WriteString(DimStyle.Render("    "+summary) + "\n")
				}
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

		// Show hint if present
		if item.hint != "" {
			b.WriteString(DimStyle.Render(indent+"  "+item.hint) + "\n")
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
