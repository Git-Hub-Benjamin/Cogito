package ui

import "github.com/charmbracelet/lipgloss"

var (
	AccentColor = lipgloss.Color("#FF6F61")
	DimColor    = lipgloss.Color("#666666")
	TextColor   = lipgloss.Color("#FFFFFF")
	BgColor     = lipgloss.Color("#1A1A2E")

	BorderStyle = lipgloss.RoundedBorder()

	BoxStyle         lipgloss.Style
	TitleStyle       lipgloss.Style
	DimStyle         lipgloss.Style
	InputPromptStyle lipgloss.Style
	ResponseStyle    lipgloss.Style
	ErrorStyle       lipgloss.Style
	SpinnerStyle     lipgloss.Style
	SelectedStyle    lipgloss.Style
	StatusBarStyle   lipgloss.Style
)

func init() {
	applyAccent()
}

// SetAccentColor updates the accent color and re-derives all dependent styles.
func SetAccentColor(hex string) {
	if hex == "" {
		return
	}
	AccentColor = lipgloss.Color(hex)
	applyAccent()
}

func applyAccent() {
	BoxStyle = lipgloss.NewStyle().
		Border(BorderStyle).
		BorderForeground(AccentColor).
		Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	DimStyle = lipgloss.NewStyle().
		Foreground(DimColor)

	InputPromptStyle = lipgloss.NewStyle().
		Foreground(AccentColor)

	ResponseStyle = lipgloss.NewStyle().
		Foreground(TextColor)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4444")).
		Bold(true)

	SpinnerStyle = lipgloss.NewStyle().
		Foreground(AccentColor)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
		Foreground(DimColor).
		Italic(true)
}

// RenderBorderTitle renders a top border line with an embedded title.
func RenderBorderTitle(title string, width int) string {
	titleRendered := " " + TitleStyle.Render(title) + " "
	titleLen := lipgloss.Width(titleRendered)

	if width < titleLen+4 {
		return titleRendered
	}

	leftBorder := lipgloss.NewStyle().Foreground(AccentColor).Render(
		BorderStyle.TopLeft + repeatStr(BorderStyle.Top, 2),
	)
	rightBorderLen := width - 4 - titleLen
	if rightBorderLen < 0 {
		rightBorderLen = 0
	}
	rightBorder := lipgloss.NewStyle().Foreground(AccentColor).Render(
		repeatStr(BorderStyle.Top, rightBorderLen) + BorderStyle.TopRight,
	)

	return leftBorder + titleRendered + rightBorder
}

func repeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
