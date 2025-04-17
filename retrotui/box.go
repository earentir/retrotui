package retrotui

import (
	"github.com/gdamore/tcell/v2"
)

// DrawBox draws a box using Unicode characters:
// if options.Doubleline is false it uses ┌, ┐, └, ┘, with horizontal (─) and vertical (│) edges, otherwise
// it uses ╔, ╗, ╚, ╝ with horizontal (═) and vertical (║) edges.
// It fills the interior with the background color (using the fill option).
// The shadow is drawn if the ShadowEnabled option is set.
// The border color is set using the borderColor parameter.
// The background color is set using the bgColor parameter.
// The box is drawn at the specified (x, y) position with the specified width (w) and height (h).
func DrawBox(s tcell.Screen, x, y, w, h int, borderColor, bgColor tcell.Color, options DrawOptions) {
	horizontalLine, verticalLine, topLeft, topRight, bottomLeft, bottomRight := '─', '│', '┌', '┐', '└', '┘'
	if options.DoubleLine {
		horizontalLine = '═'
		verticalLine = '║'
		topLeft = '╔'
		topRight = '╗'
		bottomLeft = '╚'
		bottomRight = '╝'
	}

	if w < 2 || h < 2 {
		return
	}
	borderStyle := tcell.StyleDefault.Foreground(borderColor).Background(bgColor)
	// Top & bottom edges.
	for i := 1; i < w-1; i++ {
		s.SetContent(x+i, y, horizontalLine, nil, borderStyle)
		s.SetContent(x+i, y+h-1, horizontalLine, nil, borderStyle)
	}
	// Left & right edges.
	for j := 1; j < h-1; j++ {
		s.SetContent(x, y+j, verticalLine, nil, borderStyle)
		s.SetContent(x+w-1, y+j, verticalLine, nil, borderStyle)
	}
	// Corners.
	s.SetContent(x, y, topLeft, nil, borderStyle)
	s.SetContent(x+w-1, y, topRight, nil, borderStyle)
	s.SetContent(x, y+h-1, bottomLeft, nil, borderStyle)
	s.SetContent(x+w-1, y+h-1, bottomRight, nil, borderStyle)

	// Fill interior.
	fillStyle := tcell.StyleDefault.Background(bgColor)
	fillChar := GetFillChar(options)
	for j := y + 1; j < y+h-1; j++ {
		for i := x + 1; i < x+w-1; i++ {
			s.SetContent(i, j, fillChar, nil, fillStyle)
		}
	}

	// Add shadow if enabled
	if options.ShadowEnabled {
		options.ShadowColor = tcell.ColorBlack
		options.ShadowRune = '█'
		DrawShadow(s, x, y, w, h, options)
	}
}
