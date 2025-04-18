package retrotui

import (
	"os"

	"slices"

	"github.com/gdamore/tcell/v2"
)

// DefaultKeyConfig returns a KeyConfig with reasonable default values
func DefaultKeyConfig() KeyConfig {
	return KeyConfig{
		ExitKeys: []tcell.Key{
			tcell.KeyEsc,
			tcell.KeyF3,
			tcell.KeyCtrlC,
		},
		ExitRunes: []rune{
			'q',
		},
		NavUpKey:    tcell.KeyUp,
		NavDownKey:  tcell.KeyDown,
		NavLeftKey:  tcell.KeyLeft,
		NavRightKey: tcell.KeyRight,
		SelectKey:   tcell.KeyEnter,
	}
}

// HandleBasicNavigation handles basic navigation and exit actions
// Returns:
// - true if the app should exit
// - an Action value indicating what happened
func HandleBasicNavigation(ev tcell.Event, keys KeyConfig) (bool, NavigationAction) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		// Check for exit keys
		if slices.Contains(keys.ExitKeys, e.Key()) {
			return true, NavExit
		}

		// Check for exit runes
		if slices.Contains(keys.ExitRunes, e.Rune()) {
			return true, NavExit
		}

		// Check for navigation keys
		switch e.Key() {
		case keys.NavUpKey:
			return false, NavUp
		case keys.NavDownKey:
			return false, NavDown
		case keys.NavLeftKey:
			return false, NavLeft
		case keys.NavRightKey:
			return false, NavRight
		case keys.SelectKey:
			return false, NavSelect
		}
	case *tcell.EventError:
		return true, NavExit
	}

	return false, NavNone
}

// ExitProgram cleanly exits the program
func ExitProgram(s tcell.Screen) {
	if s != nil {
		s.Fini()
	}
	os.Exit(0)
}

// InitScreen initializes a tcell screen
func InitScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = screen.Init(); err != nil {
		return nil, err
	}

	screen.EnableMouse()
	screen.Clear()

	return screen, nil
}

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

// DefaultMenuConfig returns default configuration for selection menu
func DefaultMenuConfig() MenuConfig {
	return MenuConfig{
		AppName:                "RetroTUI Application",
		CopyrightText:          "Copyright (c) 2025",
		DefaultInstructionText: "Use arrow keys to navigate, Enter to select, Esc to exit.",
		Background:             tcell.NewRGBColor(65, 70, 217), // Blue background
		TitleBarFg:             tcell.ColorYellow,
		TitleBarBg:             tcell.ColorDarkBlue,
		StatusBarFg:            tcell.ColorGreen,
		StatusBarBg:            tcell.ColorDarkBlue,
		MainSelectionBg:        tcell.ColorTeal,
		SelectionActiveBg:      tcell.ColorDarkBlue,
		SelectionNumFg:         tcell.ColorGhostWhite,
		SelectionTextFg:        tcell.ColorBlack,
		InstructionBoxFg:       tcell.ColorWhite,
		InstructionBoxBg:       tcell.ColorDarkBlue,
		ExitKeys: []tcell.Key{
			tcell.KeyEsc,
			tcell.KeyF3,
			tcell.KeyCtrlC,
		},
		ExitRunes: []rune{
			'q',
		},
		ReturnToMenuAfterSelection: true,
	}
}

// MenuState holds the current state of the selection menu
type MenuState struct {
	CurrentSelection int
	SelectedItem     string
}

// ShowSelectionMenu displays a selection menu and handles input
// Returns the selected index and text
func ShowSelectionMenu(config MenuConfig) (int, string, error) {
	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return -1, "", err
	}
	if err = screen.Init(); err != nil {
		return -1, "", err
	}
	defer screen.Fini()
	screen.EnableMouse()
	screen.Clear()

	// Initialize state
	state := MenuState{
		CurrentSelection: 0,
		SelectedItem:     "",
	}

	// Draw initial UI
	drawUI(screen, &config, state.CurrentSelection)

	// Event loop
	for {
		ev := screen.PollEvent()
		if !handleEvents(screen, ev, &config, &state) {
			break
		}
	}

	return state.CurrentSelection, state.SelectedItem, nil
}

// handleEvents processes events for the selection menu
func handleEvents(s tcell.Screen, ev tcell.Event, config *MenuConfig, state *MenuState) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		// Check for exit keys
		for _, key := range config.ExitKeys {
			if e.Key() == key {
				return false
			}
		}

		// Check for exit runes
		for _, r := range config.ExitRunes {
			if e.Rune() == r {
				return false
			}
		}

		// Arrow key navigation
		switch e.Key() {
		case tcell.KeyUp:
			if state.CurrentSelection > 0 {
				state.CurrentSelection--
			} else {
				state.CurrentSelection = len(config.MenuItems) - 1
			}
			drawUI(s, config, state.CurrentSelection)
		case tcell.KeyDown:
			if state.CurrentSelection < len(config.MenuItems)-1 {
				state.CurrentSelection++
			} else {
				state.CurrentSelection = 0
			}
			drawUI(s, config, state.CurrentSelection)
		case tcell.KeyEnter:
			// Mark item as selected
			state.SelectedItem = config.MenuItems[state.CurrentSelection].Text

			// Show message about selection
			DrawSimpleMessage(s, state.SelectedItem, tcell.ColorWhite, tcell.ColorBlue)

			// If not returning to menu, exit loop
			if !config.ReturnToMenuAfterSelection {
				return false
			}

			// Otherwise redraw and continue
			drawUI(s, config, state.CurrentSelection)
		}

	case *tcell.EventMouse:
		mouseX, mouseY := e.Position()
		buttons := e.Buttons()

		// Compute selection-dialog bounds
		width, height := s.Size()
		maxMenuLen := 0
		for _, item := range config.MenuItems {
			if len(item.Text) > maxMenuLen {
				maxMenuLen = len(item.Text)
			}
		}
		dialogWidth := maxMenuLen + 4
		dialogHeight := len(config.MenuItems) + 4
		dialogX := (width - dialogWidth) / 2
		dialogY := (height - dialogHeight) / 2

		// Check if click is within menu area
		if mouseX >= dialogX && mouseX < dialogX+dialogWidth &&
			mouseY >= dialogY && mouseY < dialogY+dialogHeight {
			menuStartY := dialogY + 2
			idx := mouseY - menuStartY
			if idx >= 0 && idx < len(config.MenuItems) {
				if state.CurrentSelection != idx {
					state.CurrentSelection = idx
					drawUI(s, config, idx)
				}
				if buttons == tcell.ButtonPrimary {
					state.SelectedItem = config.MenuItems[idx].Text
					DrawSimpleMessage(s, state.SelectedItem, tcell.ColorWhite, tcell.ColorBlue)

					if !config.ReturnToMenuAfterSelection {
						return false
					}

					drawUI(s, config, state.CurrentSelection)
				}
			}
		}
	case *tcell.EventError:
		s.Fini()
		panic(e.Error())
	}
	return true
}

// drawUI draws the complete menu UI
func drawUI(s tcell.Screen, config *MenuConfig, selected int) {
	// Background
	options := DrawOptions{
		FillPatternEnabled: true,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}
	bgStyle := tcell.StyleDefault.Background(config.Background)
	s.Fill(GetFillChar(options), bgStyle)

	// Title box
	DrawTitleBox(s, config.AppName, config.CopyrightText, config.TitleBarFg, config.TitleBarBg)

	// Selection dialog
	DrawSelectionDialog(s, config.MenuItems, selected, config.MainSelectionBg, config.SelectionActiveBg, config.SelectionNumFg, config.SelectionTextFg)

	// Instruction box
	DrawInstructionBox(s, config.MenuItems, selected, config.DefaultInstructionText, config.InstructionBoxFg, config.InstructionBoxBg)

	// Status bar
	statusText := "F3: Quit"
	for _, key := range config.ExitKeys {
		if key == tcell.KeyF3 {
			statusText = "F3: Quit"
			break
		} else if key == tcell.KeyEsc {
			statusText = "ESC: Quit"
			break
		} else if key == tcell.KeyCtrlC {
			statusText = "Ctrl+C: Quit"
			break
		}
	}
	DrawBottomBar(s, statusText, config.StatusBarFg, config.StatusBarBg)

	s.Show()
}

// DrawMessageBox shows a message box and waits for a key press
func DrawMessageBox(s tcell.Screen, message string, messageFg, messageBg, bgColor tcell.Color) {
	// Background
	options := DrawOptions{
		FillPatternEnabled: true,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}
	bgStyle := tcell.StyleDefault.Background(bgColor)
	s.Fill(GetFillChar(options), bgStyle)

	// Message box
	DrawSimpleMessage(s, message, messageFg, messageBg)
	s.Show()

	// Wait for any key press to return to the main menu
	ev := s.PollEvent()
	for {
		switch ev.(type) {
		case *tcell.EventKey, *tcell.EventMouse:
			return
		}
		ev = s.PollEvent()
	}
}
