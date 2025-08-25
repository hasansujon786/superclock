package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/styles"
)

func TimerDigit(digits string, width int, fontName FontName) string {
	var (
		output    string
		topMargin int
	)

	switch fontName {
	case BigNarrowFont, NerdFont:
		topMargin = 1
		textStyle := lipgloss.NewStyle().
			Foreground(styles.ThemeColors.Primary)

		output = RenderBigDigits(digits, fontName, textStyle)
	case DefaultFont:
		fallthrough
	default:
		textStyle := lipgloss.NewStyle().
			Foreground(styles.ThemeColors.Primary).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ThemeColors.Muted)

		output = textStyle.Render(digits)
	}

	boxStyles := lipgloss.NewStyle().
		Width(width).
		MarginTop(topMargin).
		Align(lipgloss.Center)

	return boxStyles.Render(output)
}
