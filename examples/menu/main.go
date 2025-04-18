// examples/menu/main.go — RetroTUI Selection Menu Demo
package main

import (
	"os"

	"github.com/earentir/retrotui"
	"github.com/gdamore/tcell/v2"
)

// ----------------------------------------------------------------------------
// Default instruction text
// ----------------------------------------------------------------------------
const defaultInstructionText = "Select the Base I/O Address of your audio card.\nPlease refer to your audio card manual if necessary."

// ----------------------------------------------------------------------------
// Colour scheme
// ----------------------------------------------------------------------------
var (
	AppBackground     = tcell.NewRGBColor(65, 70, 217) // Blue background
	TitleBarFg        = tcell.ColorYellow
	TitleBarBg        = tcell.ColorDarkBlue
	StatusBarFg       = tcell.ColorGreen
	StatusBarBg       = tcell.ColorDarkBlue
	MainSelectionBg   = tcell.ColorTeal
	SelectionActiveBg = tcell.ColorDarkBlue
	SelectionNumFg    = tcell.ColorGhostWhite
	SelectionTextFg   = tcell.ColorBlack
	InstructionBoxFg  = tcell.ColorWhite
	InstructionBoxBg  = tcell.ColorDarkBlue
)

// ----------------------------------------------------------------------------
// Application state
// ----------------------------------------------------------------------------
var (
	state        retrotui.UIState
	selectedItem string // holds the chosen menu item text
)

// ----------------------------------------------------------------------------
// Event handling
// ----------------------------------------------------------------------------
func handleEvents(s tcell.Screen, ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEsc, tcell.KeyF3, tcell.KeyCtrlC:
			s.Fini()
			os.Exit(0)
		case tcell.KeyUp:
			if state.CurrentSelection > 0 {
				state.CurrentSelection--
			} else {
				state.CurrentSelection = len(state.MenuItems) - 1
			}
			drawUI(s, state.CurrentSelection)
		case tcell.KeyDown:
			if state.CurrentSelection < len(state.MenuItems)-1 {
				state.CurrentSelection++
			} else {
				state.CurrentSelection = 0
			}
			drawUI(s, state.CurrentSelection)
		case tcell.KeyEnter:
			selectedItem = state.MenuItems[state.CurrentSelection].Text
			drawMessageBox(s, selectedItem)
			return false
		}

	case *tcell.EventMouse:
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

		if mouseX >= dialogX && mouseX < dialogX+dialogWidth &&
			mouseY >= dialogY && mouseY < dialogY+dialogHeight {
			menuStartY := dialogY + 2
			idx := mouseY - menuStartY
			if idx >= 0 && idx < len(state.MenuItems) {
				if state.CurrentSelection != idx {
					state.CurrentSelection = idx
					drawUI(s, idx)
				}
				if buttons == tcell.ButtonPrimary {
					selectedItem = state.MenuItems[idx].Text
					drawMessageBox(s, selectedItem)
					return false
				}
			}
		}
	case *tcell.EventError:
		s.Fini()
		panic(e.Error())
	}
	return true
}

// ----------------------------------------------------------------------------
// Drawing functions
// ----------------------------------------------------------------------------
func drawUI(s tcell.Screen, selected int) {
	// Background
	retrotui.DrawBackground(s, state.Background, true, ' ')

	// Title box
	retrotui.DrawTitleBox(s, "Hole Divers   Version 1.01", "Copyright (c) Earentir, 2025. All Rights reserved.", TitleBarFg, TitleBarBg)

	// Selection dialog
	retrotui.DrawSelectionDialog(s, state.MenuItems, selected, MainSelectionBg, SelectionActiveBg, SelectionNumFg, SelectionTextFg)

	// Instruction box
	retrotui.DrawInstructionBox(s, state.MenuItems, selected, defaultInstructionText, InstructionBoxFg, InstructionBoxBg)

	// Status bar
	retrotui.DrawBottomBar(s, "F3: Quit", StatusBarFg, StatusBarBg)

	s.Show()
}

func drawMessageBox(s tcell.Screen, message string) {
	retrotui.DrawBackground(s, state.Background, true, ' ')
	retrotui.DrawSimpleMessage(s, message, tcell.ColorWhite, tcell.ColorBlue)
	s.Show()

	// Wait for any key press to return to the main menu
	ev := s.PollEvent()
	for {
		switch ev.(type) {
		case *tcell.EventKey, *tcell.EventMouse:
			drawUI(s, state.CurrentSelection)
			return
		}
		ev = s.PollEvent()
	}
}

// ----------------------------------------------------------------------------
// Main function
// ----------------------------------------------------------------------------
func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()
	screen.EnableMouse()
	screen.Clear()

	// Initialize menu items
	state = retrotui.UIState{
		Background: AppBackground,
		MenuItems: []retrotui.MenuItem{
			{Text: "1. 10 Rounds Hell Divers Moves", Instruction: "Play 10 rounds with the standard Hell Divers move set.\nPrepare for intense tactical action!"},
			{Text: "2. 10 Rounds Random Moves", Instruction: "Play 10 rounds with completely random move generation.\nExpect the unexpected!"},
			{Text: "3. 10 Rounds Hell Divers Moves (Timer)", Instruction: "Play 10 rounds with standard Hell Divers moves and a time limit.\nSpeed is essential for this mode!"},
		},
		CurrentSelection: 0,
	}

	drawUI(screen, state.CurrentSelection)

	for {
		ev := screen.PollEvent()
		if !handleEvents(screen, ev) {
			continue
		}
	}
}
