package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/styles"
)

type ButtonEdge = [2]string

func PilledEdge() ButtonEdge {
	return ButtonEdge{"", ""}
}

type ButtonBuuilder struct {
	left   string
	right  string
	fg, bg lipgloss.Color
}

func NewButton(fg, bg lipgloss.Color) *ButtonBuuilder {
	return &ButtonBuuilder{fg: fg, bg: bg}
}

func (b *ButtonBuuilder) Edge(edges ButtonEdge) *ButtonBuuilder {
	b.left = edges[0]
	b.right = edges[1]
	return b
}

func (b *ButtonBuuilder) Render(child string) string {
	icon := lipgloss.NewStyle().
		Foreground(b.bg)

	centerItem := lipgloss.NewStyle().
		Foreground(b.fg).
		Background(b.bg)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		icon.Render(b.left),
		centerItem.Render(child),
		icon.Render(b.right),
	)
}

var ButtonStyles = NewButton(styles.ThemeColors.Black, styles.ThemeColors.Secondary).
	Edge(PilledEdge())
