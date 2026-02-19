package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type ResponseModel struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
	maxLines int
	ready    bool
}

func NewResponseModel() ResponseModel {
	return ResponseModel{}
}

func (m *ResponseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
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

func (m *ResponseModel) AppendContent(chunk string) {
	m.content += chunk
	if m.ready {
		m.viewport.SetContent(m.content)
	}
}

func (m *ResponseModel) Finalize() {}

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

// View returns the raw content (compact, no fixed-height padding).
// Use this for normal display where the box should fit the content.
func (m ResponseModel) View() string {
	return m.content
}

// PagerView returns the viewport view (fixed height, scrollable).
// Use this only when in pager mode.
func (m ResponseModel) PagerView() string {
	if !m.ready {
		return m.content
	}
	return m.viewport.View()
}

// DefaultMaxCompactLines is the default max lines before pager activates.
const DefaultMaxCompactLines = 8

// SetMaxLines sets the configurable max lines threshold.
func (m *ResponseModel) SetMaxLines(n int) {
	m.maxLines = n
}

func (m ResponseModel) maxCompactLines() int {
	if m.maxLines > 0 {
		return m.maxLines
	}
	return DefaultMaxCompactLines
}

// ContentLineCount returns the number of lines in the content.
func (m ResponseModel) ContentLineCount() int {
	if m.content == "" {
		return 0
	}
	return strings.Count(m.content, "\n") + 1
}

// Overflows returns true if the content exceeds the max compact lines.
func (m ResponseModel) Overflows() bool {
	return m.ContentLineCount() > m.maxCompactLines()
}

func (m *ResponseModel) PageDown()   { if m.ready { m.viewport.ViewDown() } }
func (m *ResponseModel) PageUp()     { if m.ready { m.viewport.ViewUp() } }
func (m *ResponseModel) ScrollDown() { if m.ready { m.viewport.LineDown(1) } }
func (m *ResponseModel) ScrollUp()   { if m.ready { m.viewport.LineUp(1) } }
func (m *ResponseModel) GotoTop()    { if m.ready { m.viewport.GotoTop() } }
func (m *ResponseModel) GotoBottom() { if m.ready { m.viewport.GotoBottom() } }

func (m ResponseModel) ScrollPercent() string {
	if !m.ready {
		return ""
	}
	if m.viewport.AtTop() {
		return "TOP"
	}
	if m.viewport.AtBottom() {
		return "END"
	}
	total := m.viewport.TotalLineCount() - m.viewport.VisibleLineCount()
	if total <= 0 {
		return "100%"
	}
	pct := m.viewport.YOffset * 100 / total
	return fmt.Sprintf("%d%%", pct)
}
