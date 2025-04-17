package retrotui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// PrintAt writes text at (x,y) using a given style.
func PrintAt(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

// PrintCentered centers text within a specified box (if boxWidth==0, uses full width).
func PrintCentered(s tcell.Screen, y, offsetX, boxWidth int, text string, style tcell.Style) {
	sw, _ := s.Size()
	if boxWidth == 0 {
		x := (sw - len(text)) / 2
		PrintAt(s, x, y, text, style)
		return
	}
	x := offsetX + (boxWidth-len(text))/2
	PrintAt(s, x, y, text, style)
}

// GetFillChar returns the rune to be used for filling areas based on options
func GetFillChar(options DrawOptions) rune {
	if options.FillPatternEnabled {
		return options.FillRune
	}
	return ' '
}

// PrintMenuTitle prints a menu title with its hotkey underlined or highlighted
func PrintMenuTitle(s tcell.Screen, x, y int, title string, hotkey rune, style tcell.Style) {
	// Find the position of the hotkey in the title
	hotkeyPos := strings.IndexRune(strings.ToLower(title), hotkey)
	if hotkeyPos >= 0 {
		// Draw the part before the hotkey
		if hotkeyPos > 0 {
			PrintAt(s, x, y, title[:hotkeyPos], style)
		}

		// Draw the hotkey with underline attribute
		hotkeyStyle := style.Underline(true)
		s.SetContent(x+hotkeyPos, y, rune(title[hotkeyPos]), nil, hotkeyStyle)

		// Draw the part after the hotkey
		if hotkeyPos < len(title)-1 {
			PrintAt(s, x+hotkeyPos+1, y, title[hotkeyPos+1:], style)
		}
	} else {
		// No hotkey found, just print the title
		PrintAt(s, x, y, title, style)
	}
}

// FillBox fills a rectangular area with the specified background color (no border).
func FillBox(s tcell.Screen, x, y, w, h int, bgColor tcell.Color, options DrawOptions) {
	fillStyle := tcell.StyleDefault.Background(bgColor)
	fillChar := GetFillChar(options)
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			s.SetContent(i, j, fillChar, nil, fillStyle)
		}
	}
}
