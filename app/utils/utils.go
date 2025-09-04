package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
)

func If[T comparable](cond bool, item, item2 T) T {
	if cond {
		return item
	}
	return item2
}

func SpaceBetween(width int, items ...string) []string {
	if len(items) == 0 {
		return []string{}
	}
	if len(items) == 1 {
		// Single item: just return it
		return []string{items[0]}
	}

	// total length of all items
	totalLen := 0
	for _, it := range items {
		totalLen += lipgloss.Width(it)
	}

	spaceLeft := width - totalLen
	if spaceLeft < 0 {
		// Items longer than width: no extra spaces
		return items
	}

	gaps := len(items) - 1
	spacePerGap := spaceLeft / gaps
	extra := spaceLeft % gaps

	parts := []string{}
	for i, it := range items {
		parts = append(parts, it)
		if i < gaps {
			space := spacePerGap
			if i < extra {
				space++
			}
			parts = append(parts, strings.Repeat(" ", space))
		}
	}

	return parts
}

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

func NotifyAppMounted() {
	content := []byte("App Started!\n")
	if err := os.WriteFile(".temp", content, 0644); err != nil {
		panic(err)
	}
}
