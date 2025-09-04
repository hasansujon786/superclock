package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/utils"
)

func TimerDigit(digits string, width int, fontName constants.FontName) string {
	var (
		output    string
		topMargin int
	)

	switch fontName {
	case constants.BigNarrowFont, constants.NerdFont:
		topMargin = 1
		textStyle := lipgloss.NewStyle().
			Foreground(styles.ThemeColors.Primary)

		output = utils.RenderBigDigits(digits, fontName, textStyle)
	case constants.DefaultFont:
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
