package retrotui

import (
	"github.com/gdamore/tcell/v2"
)

// DrawShadow draws a drop shadow using the full block character (█) in black.
// The shadow is drawn at an offset of (1,1) relative to the given rectangle.
func DrawShadow(s tcell.Screen, x, y, w, h int, options DrawOptions) {
	shadowStyle := tcell.StyleDefault.Foreground(options.ShadowColor).Background(tcell.ColorBlack)
	sw, sh := s.Size() // Get screen dimensions

	// Draw right shadow: column x+w at rows y+1 to y+h.
	for j := y + 1; j < y+h+1 && j < sh; j++ {
		s.SetContent(x+w, j, options.ShadowRune, nil, shadowStyle)
	}
	// Draw bottom shadow: row y+h at columns x+1 to x+w.
	for i := x + 1; i < x+w+1 && i < sw; i++ {
		s.SetContent(i, y+h, options.ShadowRune, nil, shadowStyle)
	}
}
