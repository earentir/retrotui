package main

import (
	"github.com/earentir/retrotui"
	"github.com/gdamore/tcell/v2"
)

// Define our state struct
type UIState struct {
	CurrentSelection int
	MenuItems        []retrotui.MenuItem
	Background       tcell.Color
	ExitRequested    bool
}

func main() {
	// Initialize tcell screen
	screen, err := retrotui.InitScreen()
	if err != nil {
		panic(err)
	}
	defer screen.Fini()

	// Set up key configuration
	keyConfig := retrotui.DefaultKeyConfig()

	// Add 'x' as an exit key
	keyConfig.ExitRunes = append(keyConfig.ExitRunes, 'x')

	// Initialize our state
	state := UIState{
		CurrentSelection: 0,
		Background:       tcell.NewRGBColor(65, 70, 217), // Blue background
		MenuItems: []retrotui.MenuItem{
			{Text: "1. 10 Rounds Hell Divers Moves", Instruction: "Play 10 rounds with the standard Hell Divers move set.\nPrepare for intense tactical action!"},
			{Text: "2. 10 Rounds Random Moves", Instruction: "Play 10 rounds with completely random move generation.\nExpect the unexpected!"},
			{Text: "3. 10 Rounds Hell Divers Moves (Timer)", Instruction: "Play 10 rounds with standard Hell Divers moves and a time limit.\nSpeed is essential for this mode!"},
		},
	}

	// Define our colors
	titleBarFg := tcell.ColorYellow
	titleBarBg := tcell.ColorDarkBlue
	statusBarFg := tcell.ColorGreen
	statusBarBg := tcell.ColorDarkBlue
	mainSelectionBg := tcell.ColorTeal
	selectionActiveBg := tcell.ColorDarkBlue
	selectionNumFg := tcell.ColorGhostWhite
	selectionTextFg := tcell.ColorBlack
	instructionBoxFg := tcell.ColorWhite
	instructionBoxBg := tcell.ColorDarkBlue

	// Draw initial UI
	drawUI(screen, &state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
		mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
		instructionBoxFg, instructionBoxBg)

	// Event loop
	for !state.ExitRequested {
		ev := screen.PollEvent()

		// Handle navigation
		shouldExit, action := retrotui.HandleBasicNavigation(ev, keyConfig)
		if shouldExit {
			break
		}

		// Process navigation actions
		switch action {
		case retrotui.NavUp:
			if state.CurrentSelection > 0 {
				state.CurrentSelection--
			} else {
				state.CurrentSelection = len(state.MenuItems) - 1
			}
			drawUI(screen, &state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
				mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
				instructionBoxFg, instructionBoxBg)

		case retrotui.NavDown:
			if state.CurrentSelection < len(state.MenuItems)-1 {
				state.CurrentSelection++
			} else {
				state.CurrentSelection = 0
			}
			drawUI(screen, &state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
				mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
				instructionBoxFg, instructionBoxBg)

		case retrotui.NavSelect:
			// Handle selection
			selectedItem := state.MenuItems[state.CurrentSelection].Text
			retrotui.DrawMessageBox(screen, selectedItem, tcell.ColorWhite, tcell.ColorBlue, state.Background)
			drawUI(screen, &state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
				mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
				instructionBoxFg, instructionBoxBg)
		}

		// Handle mouse events separately
		switch e := ev.(type) {
		case *tcell.EventMouse:
			handleMouseEvent(screen, e, &state, mainSelectionBg, selectionActiveBg,
				selectionNumFg, selectionTextFg, titleBarFg, titleBarBg, statusBarFg,
				statusBarBg, instructionBoxFg, instructionBoxBg)
		}
	}
}

// drawUI draws the complete UI
func drawUI(s tcell.Screen, state *UIState, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
	mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
	instructionBoxFg, instructionBoxBg tcell.Color) {

	// Background
	options := retrotui.DrawOptions{
		FillPatternEnabled: true,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}
	bgStyle := tcell.StyleDefault.Background(state.Background)
	s.Fill(retrotui.GetFillChar(options), bgStyle)

	// Title box
	retrotui.DrawTitleBox(s, "Hole Divers   Version 1.01", "Copyright (c) Earentir, 2025. All Rights reserved.", titleBarFg, titleBarBg)

	// Selection dialog
	retrotui.DrawSelectionDialog(s, state.MenuItems, state.CurrentSelection, mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg)

	// Instruction box
	defaultInstructionText := "Select the Base I/O Address of your audio card.\nPlease refer to your audio card manual if necessary."
	retrotui.DrawInstructionBox(s, state.MenuItems, state.CurrentSelection, defaultInstructionText, instructionBoxFg, instructionBoxBg)

	// Status bar
	retrotui.DrawBottomBar(s, "F3: Quit", statusBarFg, statusBarBg)

	s.Show()
}

// handleMouseEvent handles mouse events
func handleMouseEvent(s tcell.Screen, e *tcell.EventMouse, state *UIState,
	mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
	titleBarFg, titleBarBg, statusBarFg, statusBarBg,
	instructionBoxFg, instructionBoxBg tcell.Color) {

	mouseX, mouseY := e.Position()
	buttons := e.Buttons()

	// Compute selection-dialog bounds
	width, height := s.Size()
	maxMenuLen := 0
	for _, item := range state.MenuItems {
		if len(item.Text) > maxMenuLen {
			maxMenuLen = len(item.Text)
		}
	}
	dialogWidth := maxMenuLen + 4
	dialogHeight := len(state.MenuItems) + 4
	dialogX := (width - dialogWidth) / 2
	dialogY := (height - dialogHeight) / 2

	// Check if click is within menu area
	if mouseX >= dialogX && mouseX < dialogX+dialogWidth &&
		mouseY >= dialogY && mouseY < dialogY+dialogHeight {
		menuStartY := dialogY + 2
		idx := mouseY - menuStartY
		if idx >= 0 && idx < len(state.MenuItems) {
			if state.CurrentSelection != idx {
				state.CurrentSelection = idx
				drawUI(s, state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
					mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
					instructionBoxFg, instructionBoxBg)
			}
			if buttons == tcell.ButtonPrimary {
				selectedItem := state.MenuItems[idx].Text
				retrotui.DrawMessageBox(s, selectedItem, tcell.ColorWhite, tcell.ColorBlue, state.Background)
				drawUI(s, state, titleBarFg, titleBarBg, statusBarFg, statusBarBg,
					mainSelectionBg, selectionActiveBg, selectionNumFg, selectionTextFg,
					instructionBoxFg, instructionBoxBg)
			}
		}
	}
}
