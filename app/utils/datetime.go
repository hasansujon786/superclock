package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
)

var zeroDuration time.Duration = 0

func DurationEnded(d time.Duration) bool {
	return d == zeroDuration
}

func FormatStopwatch(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()/10) % 100 // hundredths of a second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d:%02d", h, m, s, ms)
	}
	return fmt.Sprintf("%02d:%02d:%02d", m, s, ms)
}

// formatDuration shows hh:mm:ss
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func FormatDurationHumanize(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	parts := []string{}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}

	return strings.Join(parts, ",")
}

func FormatDurationHumanizeStyled(d time.Duration) string {
	greenStyle := lipgloss.NewStyle().Foreground(styles.ThemeColors.Success)
	muted := lipgloss.NewStyle().Foreground(styles.ThemeColors.Muted)

	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	parts := []string{}
	if h > 0 {
		parts = append(parts, RenderBigDigits(fmt.Sprintf("%d", h), constants.NerdFont, greenStyle)+muted.Render("h"))
	}
	if m > 0 {
		parts = append(parts, RenderBigDigits(fmt.Sprintf("%d", m), constants.NerdFont, greenStyle)+muted.Render("m"))
	}
	if s > 0 {
		parts = append(parts, RenderBigDigits(fmt.Sprintf("%d", s), constants.NerdFont, greenStyle)+muted.Render("s"))
	}

	return strings.Join(parts, muted.Render(","))
}
