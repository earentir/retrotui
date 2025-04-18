package retrotui

import (
	"github.com/gdamore/tcell/v2"
)

// WindowState represents the current state of a window
type WindowState int

const (
	windowStateNormal WindowState = iota
	windowStateMaximized
	windowStateMinimized
)

// ControlButton represents the type of control button
type ControlButton int

const (
	buttonMinimize ControlButton = iota
	buttonMaximize
	buttonClose
)

// Window represents a resizable, movable window in the UI
type Window struct {
	Title      string
	X          int
	Y          int
	Width      int
	Height     int
	MinWidth   int
	MinHeight  int
	State      WindowState
	Visible    bool
	Active     bool
	Dragging   bool
	Resizing   bool
	LastMouseX int
	LastMouseY int
	Content    func(s tcell.Screen, x, y, width, height int) // Function to draw window content
}

// NewWindow creates a new window with default values
func NewWindow(title string, x, y, width, height int) *Window {
	return &Window{
		Title:     title,
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		MinWidth:  20,
		MinHeight: 5,
		State:     windowStateNormal,
		Visible:   true,
		Active:    true,
	}
}

// Draw renders the window on the screen
func (w *Window) Draw(s tcell.Screen, borderFg, borderBg, titleFg, titleBg, controlFg, controlBg tcell.Color) {
	if !w.Visible {
		return
	}

	// Adjust x, y, width, height according to window state
	x, y, width, height := w.GetDimensions(s)

	// Fill the window content area with the background color
	fillOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}
	// Fill the entire window including borders
	FillBox(s, x, y, width, height, borderBg, fillOptions)

	// Draw the window border
	DrawWindowBox(s, x, y, width, height, w.Title, w.Active, borderFg, borderBg, titleFg, titleBg, controlFg, controlBg)

	// Draw window content if defined
	if w.Content != nil {
		// Content area is the inner area of the window
		contentX := x + 1
		contentY := y + 1
		contentWidth := width - 2
		contentHeight := height - 2
		w.Content(s, contentX, contentY, contentWidth, contentHeight)
	}
}

// GetDimensions returns the actual dimensions of the window based on its state
func (w *Window) GetDimensions(s tcell.Screen) (x, y, width, height int) {
	switch w.State {
	case windowStateNormal:
		return w.X, w.Y, w.Width, w.Height
	case windowStateMaximized:
		screenWidth, screenHeight := s.Size()
		return 0, 1, screenWidth, screenHeight - 2 // Leave space for menu bar and status bar
	case windowStateMinimized:
		// For minimized, we still need to return something reasonable
		// In a real implementation, you might have a task bar or just show title
		return w.X, w.Y, w.Width, 1
	default:
		return w.X, w.Y, w.Width, w.Height
	}
}

// HandleEvent processes mouse events for the window
func (w *Window) HandleEvent(ev tcell.Event, windows []*Window) bool {
	if !w.Visible {
		return false
	}

	switch e := ev.(type) {
	case *tcell.EventMouse:
		mouseX, mouseY := e.Position()
		buttons := e.Buttons()

		// Get current dimensions based on window state
		x, y, width, height := w.GetDimensions(nil)

		// Check if the window should be made active (clicked anywhere in the window)
		if buttons == tcell.ButtonPrimary && mouseX >= x && mouseX < x+width && mouseY >= y && mouseY < y+height {
			// Make this window active and bring to front
			w.Active = true
			// Return true to indicate the event was handled
			// and the caller should re-order the windows
			return true
		}

		if w.Active {
			// Handle resizing via bottom-right corner
			if w.State == windowStateNormal &&
				mouseX == x+width-1 && mouseY == y+height-1 {

				if buttons == tcell.ButtonPrimary {
					if !w.Resizing {
						w.Resizing = true
						w.LastMouseX = mouseX
						w.LastMouseY = mouseY
					}
				} else if w.Resizing {
					w.Resizing = false
				}
			}

			// Handle dragging via title bar
			if w.State == windowStateNormal &&
				mouseY == y && mouseX >= x+2 && mouseX < x+2+len(w.Title)+4 {

				if buttons == tcell.ButtonPrimary {
					if !w.Dragging {
						w.Dragging = true
						w.LastMouseX = mouseX
						w.LastMouseY = mouseY
					}
				} else if w.Dragging {
					w.Dragging = false
				}
			}

			// Handle control buttons (-, +, *)
			if mouseY == y && buttons == tcell.ButtonPrimary {
				buttonWidth := 7 // Width of each control button area "═[ X ]═"

				// Minimize button
				minButtonX := x + width - 3*buttonWidth
				if mouseX >= minButtonX && mouseX < minButtonX+buttonWidth {
					w.State = windowStateMinimized
					return true
				}

				// Maximize/restore button
				maxButtonX := x + width - 2*buttonWidth
				if mouseX >= maxButtonX && mouseX < maxButtonX+buttonWidth {
					if w.State == windowStateMaximized {
						w.State = windowStateNormal
					} else {
						w.State = windowStateMaximized
					}
					return true
				}

				// Close button
				closeButtonX := x + width - buttonWidth
				if mouseX >= closeButtonX && mouseX < closeButtonX+buttonWidth {
					w.Visible = false
					return true
				}
			}

			// Handle ongoing dragging or resizing
			if buttons == tcell.ButtonPrimary {
				if w.Dragging {
					deltaX := mouseX - w.LastMouseX
					deltaY := mouseY - w.LastMouseY
					w.X += deltaX
					w.Y += deltaY
					w.LastMouseX = mouseX
					w.LastMouseY = mouseY
					return true
				} else if w.Resizing {
					deltaX := mouseX - w.LastMouseX
					deltaY := mouseY - w.LastMouseY

					// Ensure we don't resize below minimum dimensions
					newWidth := w.Width + deltaX
					newHeight := w.Height + deltaY

					if newWidth >= w.MinWidth {
						w.Width = newWidth
						w.LastMouseX = mouseX
					}

					if newHeight >= w.MinHeight {
						w.Height = newHeight
						w.LastMouseY = mouseY
					}

					return true
				}
			}
		}
	}

	return false
}

// DrawWindowBox draws a window with title and control buttons
func DrawWindowBox(s tcell.Screen, x, y, width, height int, title string, active bool,
	borderFg, borderBg, titleFg, titleBg, controlFg, controlBg tcell.Color) {

	if width < 10 || height < 3 {
		return // Too small to draw properly
	}

	// Select colors based on active state
	currentBorderFg := borderFg
	currentBorderBg := borderBg
	currentTitleFg := titleFg
	currentTitleBg := titleBg

	if !active {
		// Use less saturated colors for inactive windows
		currentBorderFg = tcell.ColorDarkGray
		currentTitleFg = tcell.ColorGray
	}

	// Fill the inner content area with background color
	fillStyle := tcell.StyleDefault.Background(currentBorderBg)
	for j := y + 1; j < y+height-1; j++ {
		for i := x + 1; i < x+width-1; i++ {
			s.SetContent(i, j, ' ', nil, fillStyle)
		}
	}

	// Draw the control buttons
	drawControlButtons(s, x, y, width, currentBorderFg, currentBorderBg, controlFg, controlBg)

	// Draw the title
	titleSt := tcell.StyleDefault.Foreground(currentTitleFg).Background(currentTitleBg)
	titleWithBrackets := "[ " + title + " ]"
	PrintAt(s, x+2, y, titleWithBrackets, titleSt)

	// Draw top border (except where the title and controls are)
	borderSt := tcell.StyleDefault.Foreground(currentBorderFg).Background(currentBorderBg)
	s.SetContent(x, y, '╔', nil, borderSt)

	// Top border before title
	for i := 1; i < 2; i++ {
		s.SetContent(x+i, y, '═', nil, borderSt)
	}

	// Top border after title
	buttonStart := width - 21 // Start of minimize/maximize/close buttons
	for i := 2 + len(titleWithBrackets); i < buttonStart; i++ {
		s.SetContent(x+i, y, '═', nil, borderSt)
	}

	// Right top corner
	s.SetContent(x+width-1, y, '╗', nil, borderSt)

	// Left and right borders
	for j := 1; j < height-1; j++ {
		s.SetContent(x, y+j, '║', nil, borderSt)
		s.SetContent(x+width-1, y+j, '║', nil, borderSt)
	}

	// Bottom border with corners
	s.SetContent(x, y+height-1, '╚', nil, borderSt)
	for i := 1; i < width-1; i++ {
		s.SetContent(x+i, y+height-1, '═', nil, borderSt)
	}
	s.SetContent(x+width-1, y+height-1, '╝', nil, borderSt)

	// Make the bottom-right corner a special character for resizing
	s.SetContent(x+width-1, y+height-1, '╬', nil, borderSt)
}

// drawControlButtons draws the minimize, maximize, and close buttons
func drawControlButtons(s tcell.Screen, x, y, width int,
	borderFg, borderBg, controlFg, controlBg tcell.Color) {

	borderSt := tcell.StyleDefault.Foreground(borderFg).Background(borderBg)
	controlSt := tcell.StyleDefault.Foreground(controlFg).Background(controlBg)

	// Minimize button
	minX := width - 21
	for i := 0; i < 2; i++ {
		s.SetContent(x+minX+i, y, '═', nil, borderSt)
	}
	s.SetContent(x+minX+2, y, '[', nil, borderSt)
	s.SetContent(x+minX+3, y, ' ', nil, borderSt)
	s.SetContent(x+minX+4, y, '-', nil, controlSt)
	s.SetContent(x+minX+5, y, ' ', nil, borderSt)
	s.SetContent(x+minX+6, y, ']', nil, borderSt)

	// Maximize button
	maxX := width - 14
	for i := 0; i < 2; i++ {
		s.SetContent(x+maxX+i, y, '═', nil, borderSt)
	}
	s.SetContent(x+maxX+2, y, '[', nil, borderSt)
	s.SetContent(x+maxX+3, y, ' ', nil, borderSt)
	s.SetContent(x+maxX+4, y, '+', nil, controlSt)
	s.SetContent(x+maxX+5, y, ' ', nil, borderSt)
	s.SetContent(x+maxX+6, y, ']', nil, borderSt)

	// Close button
	closeX := width - 7
	for i := 0; i < 2; i++ {
		s.SetContent(x+closeX+i, y, '═', nil, borderSt)
	}
	s.SetContent(x+closeX+2, y, '[', nil, borderSt)
	s.SetContent(x+closeX+3, y, ' ', nil, borderSt)
	s.SetContent(x+closeX+4, y, '*', nil, controlSt)
	s.SetContent(x+closeX+5, y, ' ', nil, borderSt)
	s.SetContent(x+closeX+6, y, ']', nil, borderSt)
}

// ManageWindows handles window z-order and event routing to appropriate windows
func ManageWindows(s tcell.Screen, windows []*Window, ev tcell.Event,
	borderFg, borderBg, titleFg, titleBg, controlFg, controlBg tcell.Color) bool {

	// Start from the top window (last in the array) and work backwards
	for i := len(windows) - 1; i >= 0; i-- {
		window := windows[i]

		// Try to handle the event with this window
		if window.HandleEvent(ev, windows) {
			// If the event was handled, move this window to the top of the z-order
			if i < len(windows)-1 {
				// Remove the window from its current position
				windowToMove := windows[i]
				windows = append(windows[:i], windows[i+1:]...)
				// Add it to the end (top of z-order)
				windows = append(windows, windowToMove)
			}

			// Redraw all windows
			DrawWindows(s, windows, borderFg, borderBg, titleFg, titleBg, controlFg, controlBg)
			return true
		}
	}

	return false
}

// DrawWindows draws all visible windows in z-order
func DrawWindows(s tcell.Screen, windows []*Window,
	borderFg, borderBg, titleFg, titleBg, controlFg, controlBg tcell.Color) {

	// Draw from bottom to top
	for _, window := range windows {
		if window.Visible {
			window.Draw(s, borderFg, borderBg, titleFg, titleBg, controlFg, controlBg)
		}
	}
}
