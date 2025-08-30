package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/utils"
)

type Digits map[rune][]string
type FontName int

const (
	DefaultFont FontName = iota
	BigNarrowFont
	NerdFont
)

var NERD_FONTS = Digits{
	'0': {"󰬹"},
	'1': {"󰬺"},
	'2': {"󰬻"},
	'3': {"󰬼"},
	'4': {"󰬽"},
	'5': {"󰬾"},
	'6': {"󰬿"},
	'7': {"󰭀"},
	'8': {"󰭁"},
	'9': {"󰭂"},
	':': {":"}, // colon
}

var BIG_NARROW_FONT = Digits{
	'0': {
		"███",
		"█ █",
		"█ █",
		"█ █",
		"███",
	},
	'1': {
		"▄█ ",
		" █ ",
		" █ ",
		" █ ",
		"▄█▄",
	},
	'2': {
		"███",
		"  █",
		"███",
		"█  ",
		"███",
	},
	'3': {
		"███",
		"  █",
		"███",
		"  █",
		"███",
	},
	'4': {
		"█ █",
		"█ █",
		"███",
		"  █",
		"  █",
	},
	'5': {
		"███",
		"█  ",
		"███",
		"  █",
		"███",
	},
	'6': {
		"███",
		"█  ",
		"███",
		"█ █",
		"███",
	},
	'7': {
		"███",
		"  █",
		"  █",
		"  █",
		"  █",
	},
	'8': {
		"███",
		"█ █",
		"███",
		"█ █",
		"███",
	},
	'9': {
		"███",
		"█ █",
		"███",
		"  █",
		"███",
	},
	':': { // optional for stopwatch display
		"   ",
		" █ ",
		"   ",
		" █ ",
		"   ",
	},
}

func getDigits(fontName FontName) Digits {
	switch fontName {
	case BigNarrowFont:
		return BIG_NARROW_FONT
	case DefaultFont:
		fallthrough
	default:
		return NERD_FONTS
	}
}

// func RenderBigDigitsOld(s string, style lipgloss.Style) string {
// 	lines := make([]string, 5)
// 	// width := (len(s) * 3) + (len(s) - 1)
//
// 	for row := range 5 {
// 		parts := []string{}
// 		// parts = append(parts, "xxx")
// 		for _, ch := range s {
// 			if glyph, ok := BIG_NARROW_FONT[ch]; ok {
// 				parts = append(parts, style.Render(glyph[row]))
// 				// add space after each digit
// 				parts = append(parts, style.Render(" "))
// 			}
// 		}
// 		lines[row] = lipgloss.JoinHorizontal(lipgloss.Left, parts...)
// 	}
//
// 	return lipgloss.JoinVertical(lipgloss.Left, lines...)
// }

// Render function using Lip Gloss
func RenderBigDigits(s string, fontName FontName, style lipgloss.Style) string {
	digits := getDigits(fontName)

	lines := []string{}
	maxRows := 0
	// Determine max number of rows for this font
	for _, ch := range s {
		if glyph, ok := digits[ch]; ok {
			if len(glyph) > maxRows {
				maxRows = len(glyph)
			}
		}
	}

	// Build lines
	for row := 0; row < maxRows; row++ {
		parts := []string{}
		for _, ch := range s {
			if glyph, ok := digits[ch]; ok {
				line := ""
				if row < len(glyph) {
					line = glyph[row]
				}

				// Don't Add space if char is : and NF
				isColonAndNF := line == ":" && fontName == NerdFont
				parts = append(parts, style.Render(line)+utils.If(isColonAndNF, "", " "))
			}
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, parts...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
