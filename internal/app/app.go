package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/benji/cogito/internal/config"
	shellctx "github.com/benji/cogito/internal/context"
	"github.com/benji/cogito/internal/provider"
	"github.com/benji/cogito/internal/ui"
)

var commands = []string{"/settings", "/help", "/clear"}

type Model struct {
	state    AppState
	config   config.Config
	provider *provider.OpenAIProvider

	input    ui.InputModel
	response ui.ResponseModel
	settings ui.SettingsModel
	spinner  spinner.Model

	width  int
	height int

	streamCh    <-chan string
	streamErrCh <-chan error
	cancelFunc  context.CancelFunc

	err      error
	hasError bool

	// topInline mode: top position without clear screen.
	// Box is half-height and scroll is locked to keep render size fixed.
	topInline bool
}

func NewModel(cfg config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.SpinnerStyle

	ui.SetAccentColor(cfg.Theme.AccentColor)
	p := provider.NewOpenAI(cfg.APIKey(), cfg.DefaultModel, cfg.BaseURL)

	return Model{
		state:     StateInput,
		config:    cfg,
		provider:  p,
		input:     ui.NewInputModel(),
		response:  ui.NewResponseModel(),
		spinner:   s,
		topInline: !cfg.ClearScreen && cfg.Position == "top",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentWidth := m.width - 6 // borders + padding
		if contentWidth < 20 {
			contentWidth = 20
		}
		m.input.SetWidth(contentWidth)

		// Determine usable height for the box
		usableHeight := m.height
		if m.topInline {
			// Cap to half the terminal so old content stays visible below
			usableHeight = m.height / 2
		}

		// Reserve space for header, input, status bar
		responseHeight := usableHeight - 8
		if responseHeight < 3 {
			responseHeight = 3
		}
		m.response.SetSize(contentWidth, responseHeight)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case streamChunkMsg:
		m.response.AppendContent(msg.chunk)
		return m, listenForChunks(m.streamCh, m.streamErrCh)

	case streamDoneMsg:
		m.state = StateInput
		m.streamCh = nil
		m.streamErrCh = nil
		m.cancelFunc = nil
		return m, m.input.Focus()

	case streamErrMsg:
		m.state = StateInput
		m.err = msg.err
		m.hasError = true
		m.streamCh = nil
		m.streamErrCh = nil
		m.cancelFunc = nil
		return m, m.input.Focus()

	case ui.SettingsSavedMsg:
		if m.config.APIKeys == nil {
			m.config.APIKeys = make(map[string]string)
		}
		m.config.APIKeys["openai"] = msg.APIKey
		m.config.BaseURL = msg.BaseURL
		m.config.DefaultModel = msg.DefaultModel
		m.config.ClearScreen = msg.ClearScreen
		m.config.Position = msg.Position
		m.config.Theme.AccentColor = msg.AccentColor
		ui.SetAccentColor(msg.AccentColor)
		m.provider = provider.NewOpenAI(msg.APIKey, msg.DefaultModel, msg.BaseURL)
		m.state = StateInput
		_ = m.config.Save()
		m.hasError = false
		return m, m.input.Focus()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m.updateSubmodels(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateInput:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+k":
			return m, m.input.Focus()
		case "tab":
			return m.handleTabComplete()
		case "enter":
			return m.handleSubmit()
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case StateStreaming:
		switch msg.String() {
		case "esc":
			if m.cancelFunc != nil {
				m.cancelFunc()
			}
			m.state = StateInput
			return m, m.input.Focus()
		case "ctrl+c":
			return m, tea.Quit
		}
		// In topInline mode, skip scroll keys to keep render height fixed
		if m.topInline {
			return m, nil
		}
		var cmd tea.Cmd
		m.response, cmd = m.response.Update(msg)
		return m, cmd

	case StateSettings:
		switch msg.String() {
		case "esc":
			m.state = StateInput
			return m, m.input.Focus()
		case "ctrl+c":
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleTabComplete() (tea.Model, tea.Cmd) {
	val := m.input.Value()
	if val == "" || val[0] != '/' {
		return m, nil
	}

	prefix := strings.ToLower(val)
	var matches []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, cmd)
		}
	}

	if len(matches) == 1 {
		m.input.SetValue(matches[0])
	} else if len(matches) > 1 {
		// Complete to longest common prefix
		common := matches[0]
		for _, match := range matches[1:] {
			common = commonPrefix(common, match)
		}
		if len(common) > len(val) {
			m.input.SetValue(common)
		}
	}

	return m, nil
}

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}

func (m Model) handleSubmit() (tea.Model, tea.Cmd) {
	query := strings.TrimSpace(m.input.Value())
	if query == "" {
		return m, nil
	}

	// Handle commands
	switch {
	case query == "/settings":
		m.settings = ui.NewSettingsModel(
			m.config.APIKey(), m.config.BaseURL, m.config.DefaultModel,
			m.config.ClearScreen, m.config.Position, m.config.Theme.AccentColor,
		)
		contentWidth := m.width - 6
		if contentWidth > 0 {
			m.settings.SetWidth(contentWidth)
		}
		m.state = StateSettings
		m.input.SetValue("")
		m.input.Blur()
		return m, nil

	case query == "/clear":
		m.response.Clear()
		m.input.SetValue("")
		m.hasError = false
		return m, nil

	case query == "/help":
		m.response.Clear()
		m.response.AppendContent(helpText())
		m.input.SetValue("")
		return m, nil
	}

	// Check for API key
	if m.config.APIKey() == "" {
		m.err = fmt.Errorf("no API key set — run /settings or set OPENAI_API_KEY")
		m.hasError = true
		m.input.SetValue("")
		return m, nil
	}

	// Start streaming
	m.response.Clear()
	m.hasError = false
	m.state = StateStreaming
	m.input.SetValue("")
	m.input.Blur()

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel

	ch := make(chan string, 64)
	errCh := make(chan error, 1)
	m.streamCh = ch
	m.streamErrCh = errCh

	messages := []provider.ChatMessage{
		{Role: provider.RoleSystem, Content: buildSystemMsg(m.config.Context.IncludeCWD)},
		{Role: provider.RoleUser, Content: query},
	}

	go func() {
		errCh <- m.provider.StreamChat(ctx, messages, ch)
	}()

	return m, listenForChunks(ch, errCh)
}

func (m Model) updateSubmodels(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case StateInput:
		m.input, cmd = m.input.Update(msg)
	case StateStreaming:
		m.response, cmd = m.response.Update(msg)
	case StateSettings:
		m.settings, cmd = m.settings.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	contentWidth := m.width - 6
	if contentWidth < 20 {
		contentWidth = 20
	}

	title := ui.RenderHeader(m.config.DefaultModel)
	topBorder := ui.RenderBorderTitle(title, m.width)

	var content string
	switch m.state {
	case StateSettings:
		content = m.settings.View()
	default:
		content = m.buildMainView(contentWidth)
	}

	body := lipgloss.NewStyle().
		Width(contentWidth).
		Padding(0, 2).
		Render(content)

	// Build bottom border
	bottomBorder := lipgloss.NewStyle().Foreground(ui.AccentColor).Render(
		ui.BorderStyle.BottomLeft +
			repeatStr(ui.BorderStyle.Bottom, m.width-2) +
			ui.BorderStyle.BottomRight,
	)

	// Side borders
	lines := strings.Split(body, "\n")
	var bordered strings.Builder
	bordered.WriteString(topBorder + "\n")
	leftBar := lipgloss.NewStyle().Foreground(ui.AccentColor).Render(ui.BorderStyle.Left)
	rightBar := lipgloss.NewStyle().Foreground(ui.AccentColor).Render(ui.BorderStyle.Right)
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		pad := contentWidth + 4 - lineWidth
		if pad < 0 {
			pad = 0
		}
		bordered.WriteString(leftBar + line + strings.Repeat(" ", pad) + rightBar + "\n")
	}
	bordered.WriteString(bottomBorder)

	result := bordered.String()

	// Push to bottom of screen when using alt screen + bottom position
	if m.config.ClearScreen && m.config.Position == "bottom" {
		renderedLines := strings.Count(result, "\n") + 1
		padding := m.height - renderedLines
		if padding > 0 {
			result = strings.Repeat("\n", padding) + result
		}
	}

	return result
}

func (m Model) buildMainView(width int) string {
	var parts []string

	// Response area
	if m.response.Content() != "" || m.state == StateStreaming {
		responseView := m.response.View()
		if m.state == StateStreaming {
			responseView += m.spinner.View()
		}
		parts = append(parts, responseView)
	}

	// Error
	if m.hasError && m.err != nil {
		parts = append(parts, ui.ErrorStyle.Render("Error: "+m.err.Error()))
	}

	// Input
	parts = append(parts, m.input.View())

	// Status bar
	status := m.statusBar()
	if status != "" {
		parts = append(parts, ui.StatusBarStyle.Render(status))
	}

	return strings.Join(parts, "\n\n")
}

func (m Model) statusBar() string {
	switch m.state {
	case StateStreaming:
		return "Streaming... (esc to cancel)"
	default:
		return "/help commands • /settings configure • esc quit"
	}
}

func repeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func helpText() string {
	return `Commands:
  /settings   - Configure API key, model, and preferences
  /clear      - Clear response
  /help       - Show this help

Shortcuts:
  Enter       - Submit query
  Tab         - Autocomplete commands
  Ctrl+K      - Focus input
  Esc         - Quit (or cancel streaming)`
}

func buildSystemMsg(includeCWD bool) string {
	return shellctx.BuildSystemMessage(includeCWD)
}
