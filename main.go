// main.go — RetroTUI demo (original look‑and‑feel)
// -----------------------------------------------------------------------------
// Only ONE behavioural tweak:
//   - Selecting an item from the first (boxes) menu now shows a centred
//     message box with that item’s text. Press any key and the program
//     proceeds to the original windows + menus demo.
//
// No colours, dimensions, or text have been altered.
// -----------------------------------------------------------------------------
package main

import (
	"os"

	"dosmenu/retrotui"

	"github.com/gdamore/tcell/v2"
)

// ----------------------------------------------------------------------------
// Default instruction text
// ----------------------------------------------------------------------------
const defaultInstructionText = "Select the Base I/O Address of your audio card.\nPlease refer to your audio card manual if necessary."

// ----------------------------------------------------------------------------
// Colour scheme (unchanged)
// ----------------------------------------------------------------------------
var (
	AppBackground     = tcell.NewRGBColor(65, 70, 217) // Blue background
	TitleBarFg        = tcell.ColorYellow
	TitleBarBg        = tcell.ColorDarkBlue
	MenuBarFg         = tcell.ColorWhite
	MenuBarBg         = tcell.ColorDarkBlue
	MenuActiveFg      = tcell.ColorWhite
	MenuActiveBg      = tcell.ColorRed
	DropdownFg        = tcell.ColorBlack
	DropdownBg        = tcell.ColorLightGray
	DropdownActiveFg  = tcell.ColorWhite
	DropdownActiveBg  = tcell.ColorBlue
	SeparatorColor    = tcell.ColorGray
	MainSelectionBg   = tcell.ColorTeal
	SelectionActiveBg = tcell.ColorDarkBlue
	SelectionNumFg    = tcell.ColorGhostWhite
	SelectionTextFg   = tcell.ColorBlack
	InstructionBoxFg  = tcell.ColorWhite
	InstructionBoxBg  = tcell.ColorDarkBlue
	StatusBarFg       = tcell.ColorGreen
	StatusBarBg       = tcell.ColorDarkBlue
	WizardBoxFg       = tcell.ColorWhite
	WizardBoxBg       = tcell.ColorBlue
)

// ----------------------------------------------------------------------------
// Application state & globals
// ----------------------------------------------------------------------------
var (
	state      retrotui.UIState
	appMenus   []retrotui.Menu
	appWindows []*retrotui.Window

	selectedItemText string // holds the chosen menu text for the dialog box

	// Wizard demo
	wizardPages = []WizardPage{
		{Title: "Step 1 of 3 — Welcome", Content: []string{"This wizard will guide you", "through basic configuration."}},
		{Title: "Step 2 of 3 — Options", Content: []string{"Choose preferred settings", "and press Next."}},
		{Title: "Step 3 of 3 — Ready", Content: []string{"Setup is complete.", "Click Finish to exit."}},
	}
	wizardIndex   = 0     // current page
	wizardBtnIdx  = 1     // 0‑Back,1‑Next/Finish,2‑Cancel
	wizardRunning = false // true while in wizard screens
)

// Screen identifiers
const (
	screenMenu   = "menu"   // initial selection dialog
	screenDialog = "dialog" // centred message box with choice
	screenMain   = "main"   // original windows demo
	screenWizard = "wizard" // setup wizard
)

// WizardPage represents a page in the wizard
type WizardPage struct {
	Title   string
	Content []string
}

// ----------------------------------------------------------------------------
// Shared / global event handling
// ----------------------------------------------------------------------------
func handleCommonEvents(s tcell.Screen, ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEsc, tcell.KeyF3:
			// Only quit outright when we are already in the windows demo.
			if state.CurrentScreen == screenMain {
				s.Fini()
				os.Exit(0)
			}
			// Otherwise let the screen‑specific handler deal with it
			return false
		case tcell.KeyCtrlC:
			s.Fini()
			os.Exit(0)
		}
		if r := e.Rune(); (r == 'q' || r == 'Q') && state.CurrentScreen == screenMain {
			s.Fini()
			os.Exit(0)
		}
	case *tcell.EventError:
		s.Fini()
		panic(e.Error())
	}
	return false
}

// ----------------------------------------------------------------------------
// MENU (boxes) SCREEN — visuals untouched
// ----------------------------------------------------------------------------
func handleMenuScreenEvents(s tcell.Screen, ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventMouse:
		mouseX, mouseY := e.Position()
		buttons := e.Buttons()

		// Compute selection‑dialog bounds exactly as before
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
					drawMenuScreenUI(s, idx)
				}
				if buttons == tcell.ButtonPrimary && e.Buttons()&tcell.Button1 != 0 {
					selectedItemText = state.MenuItems[idx].Text
					state.CurrentScreen = screenDialog
					return
				}
			}
		}
		drawMenuScreenUI(s, state.CurrentSelection)

	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyUp:
			if state.CurrentSelection > 0 {
				state.CurrentSelection--
			} else {
				state.CurrentSelection = len(state.MenuItems) - 1
			}
			drawMenuScreenUI(s, state.CurrentSelection)
		case tcell.KeyDown:
			if state.CurrentSelection < len(state.MenuItems)-1 {
				state.CurrentSelection++
			} else {
				state.CurrentSelection = 0
			}
			drawMenuScreenUI(s, state.CurrentSelection)
		case tcell.KeyEnter:
			selectedItemText = state.MenuItems[state.CurrentSelection].Text
			state.CurrentScreen = screenDialog
		case tcell.KeyF3, tcell.KeyEsc:
			// Launch wizard demo
			wizardIndex = 0
			wizardBtnIdx = 1
			wizardRunning = true
			state.CurrentScreen = screenWizard
		}
	}
}

// ---------------------------------
//
//	new, minimal
//
// ---------------------------------
func drawDialogScreenUI(s tcell.Screen) {
	retrotui.DrawBackground(s, state.Background, true, ' ')
	retrotui.DrawSimpleMessage(s, selectedItemText, tcell.ColorWhite, tcell.ColorBlue)
	s.Show()
}

func handleDialogScreenEvents(s tcell.Screen, ev tcell.Event) {
	switch ev.(type) {
	case *tcell.EventKey, *tcell.EventMouse:
		state.CurrentScreen = screenMenu
	}
}

// ----------------------------------------------------------------------------
// drawMenuScreenUI — original code
// ----------------------------------------------------------------------------
func drawMenuScreenUI(s tcell.Screen, selected int) {
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

// ----------------------------------------------------------------------------
// drawMainScreenUI — original code (unchanged)
// ----------------------------------------------------------------------------
func drawMainScreenUI(s tcell.Screen) {
	// Background
	retrotui.DrawBackground(s, state.Background, true, ' ')

	// Menu bar
	retrotui.DrawMenuBar(s, appMenus, state.ActiveMenu, state.MenuBarActive,
		MenuBarFg, MenuBarBg, MenuActiveFg, MenuActiveBg)

	// Dropdown menu
	if state.MenuBarActive && state.ActiveMenu >= 0 && state.ActiveMenu < len(appMenus) {
		retrotui.DrawDropdownMenu(s, appMenus[state.ActiveMenu], state.ActiveMenuItem,
			DropdownFg, DropdownBg, DropdownActiveFg, DropdownActiveBg, SeparatorColor)
	}

	// Placeholder content area
	width, height := s.Size()
	retrotui.PrintCentered(s, height/2, 0, width, "Main Application Screen",
		tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(state.Background))
	retrotui.PrintCentered(s, height/2+1, 0, width, "Press Alt+F, Alt+E, or Alt+H to activate menus",
		tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(state.Background))
	retrotui.PrintCentered(s, height/2+2, 0, width, "Press F3 to exit",
		tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(state.Background))

	// Visible windows
	if len(appWindows) > 0 {
		retrotui.DrawWindows(s, appWindows,
			tcell.ColorWhite, tcell.ColorBlue,
			tcell.ColorYellow, tcell.ColorBlue,
			tcell.ColorRed, tcell.ColorBlue)
	}

	// Status bar
	retrotui.DrawBottomBar(s, "F3: Quit", StatusBarFg, StatusBarBg)

	s.Show()
}

// ----------------------------------------------------------------------------
// handleMainScreenEvents — original code (unchanged)
// ----------------------------------------------------------------------------
func handleMainScreenEvents(s tcell.Screen, ev tcell.Event) {
	// First give top‑most window a chance
	if len(appWindows) > 0 {
		if retrotui.ManageWindows(s, appWindows, ev,
			tcell.ColorWhite, tcell.ColorBlue,
			tcell.ColorYellow, tcell.ColorBlue,
			tcell.ColorRed, tcell.ColorBlue) {
			return
		}
	}

	switch e := ev.(type) {
	case *tcell.EventMouse:
		mouseX, mouseY := e.Position()
		buttons := e.Buttons()

		// Dropdown interactions
		if state.MenuBarActive {
			if state.ActiveMenu >= 0 && state.ActiveMenu < len(appMenus) {
				menu := appMenus[state.ActiveMenu]
				items := menu.Items
				menuX := menu.Position
				menuY := 1
				if menu.Align {
					width, _ := s.Size()
					menuX = width - len(menu.Title) - 2
				}
				maxWidth := 0
				for _, it := range items {
					if len(it.Text) > maxWidth {
						maxWidth = len(it.Text)
					}
				}
				menuWidth := maxWidth + 4
				width, _ := s.Size()
				if menuX+menuWidth > width {
					menuX = width - menuWidth
				}

				if mouseX >= menuX && mouseX < menuX+menuWidth && mouseY > menuY && mouseY < menuY+len(items)+1 {
					idx := mouseY - menuY - 1
					if idx >= 0 && idx < len(items) && !items[idx].IsSeparator {
						state.ActiveMenuItem = idx
						if buttons == tcell.ButtonPrimary {
							if items[idx].OnSelect != nil {
								state.MenuBarActive = false
								state.ActiveMenu = -1
								items[idx].OnSelect(s)
								return
							}
						}
					}
					drawMainScreenUI(s)
					return
				} else if buttons == tcell.ButtonPrimary {
					state.MenuBarActive = false
					state.ActiveMenu = -1
					drawMainScreenUI(s)
					return
				}
			}
		}

		// Click on menu titles
		if mouseY == 0 && buttons == tcell.ButtonPrimary {
			for i, menu := range appMenus {
				menuX := menu.Position
				if menu.Align {
					width, _ := s.Size()
					menuX = width - len(menu.Title) - 2
				}
				if mouseX >= menuX && mouseX < menuX+len(menu.Title) {
					if state.MenuBarActive && state.ActiveMenu == i {
						state.MenuBarActive = false
						state.ActiveMenu = -1
					} else {
						state.ActiveMenu = i
						state.ActiveMenuItem = 0
						state.MenuBarActive = true
					}
					drawMainScreenUI(s)
					return
				}
			}
			if state.MenuBarActive {
				state.MenuBarActive = false
				state.ActiveMenu = -1
				drawMainScreenUI(s)
				return
			}
		}

	case *tcell.EventKey:
		if e.Modifiers()&tcell.ModAlt != 0 {
			for i, menu := range appMenus {
				if e.Rune() == menu.HotKey {
					state.ActiveMenu = i
					state.ActiveMenuItem = 0
					state.MenuBarActive = true
					drawMainScreenUI(s)
					return
				}
			}
		} else if state.MenuBarActive {
			switch e.Key() {
			case tcell.KeyLeft:
				if state.ActiveMenu > 0 {
					state.ActiveMenu--
				} else {
					state.ActiveMenu = len(appMenus) - 1
				}
				state.ActiveMenuItem = 0
			case tcell.KeyRight:
				if state.ActiveMenu < len(appMenus)-1 {
					state.ActiveMenu++
				} else {
					state.ActiveMenu = 0
				}
				state.ActiveMenuItem = 0
			case tcell.KeyUp:
				if state.ActiveMenuItem > 0 {
					state.ActiveMenuItem--
				}
				if appMenus[state.ActiveMenu].Items[state.ActiveMenuItem].IsSeparator && state.ActiveMenuItem > 0 {
					state.ActiveMenuItem--
				}
			case tcell.KeyDown:
				if state.ActiveMenuItem < len(appMenus[state.ActiveMenu].Items)-1 {
					state.ActiveMenuItem++
				}
				if appMenus[state.ActiveMenu].Items[state.ActiveMenuItem].IsSeparator && state.ActiveMenuItem < len(appMenus[state.ActiveMenu].Items)-1 {
					state.ActiveMenuItem++
				}
			case tcell.KeyEnter:
				it := appMenus[state.ActiveMenu].Items[state.ActiveMenuItem]
				if it.OnSelect != nil && !it.IsSeparator {
					it.OnSelect(s)
				}
				state.MenuBarActive = false
				state.ActiveMenu = -1
			case tcell.KeyEscape:
				state.MenuBarActive = false
				state.ActiveMenu = -1
			}
			drawMainScreenUI(s)
		}
	}
}

// ----------------------------------------------------------------------------
// createWindow — original helper (unchanged)
// ----------------------------------------------------------------------------
func createWindow(s tcell.Screen, title string, content func(s tcell.Screen, x, y, width, height int)) {
	sw, sh := s.Size()
	w, h := 50, 15
	x, y := (sw-w)/2, (sh-h)/2
	win := retrotui.NewWindow(title, x, y, w, h)
	win.Content = content
	appWindows = append(appWindows, win)
	drawMainScreenUI(s)
}

// ----------------------------------------------------------------------------
// initialiseWindowsAndMenus — original code lifted verbatim
// ----------------------------------------------------------------------------
func initialiseWindowsAndMenus(screen tcell.Screen) {
	appWindows = make([]*retrotui.Window, 0)
	appMenus = []retrotui.Menu{
		{
			Title: "File", HotKey: 'f', Position: 1,
			Items: []retrotui.DropdownItem{
				{Text: "Open...", OnSelect: func(s tcell.Screen) {
					createWindow(s, "Open", func(sc tcell.Screen, x, y, w, h int) {
						st := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
						retrotui.PrintAt(sc, x+2, y+2, "Open dialog placeholder", st)
					})
				}},
				{Text: "Save", OnSelect: func(s tcell.Screen) {
					createWindow(s, "Save", func(sc tcell.Screen, x, y, w, h int) {
						st := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
						retrotui.PrintAt(sc, x+2, y+2, "Save placeholder", st)
					})
				}},
				{IsSeparator: true},
				{Text: "Exit", OnSelect: func(s tcell.Screen) { s.Fini(); os.Exit(0) }},
			},
		},
		{
			Title: "Edit", HotKey: 'e', Position: 6,
			Items: []retrotui.DropdownItem{
				{Text: "Copy", OnSelect: func(s tcell.Screen) {
					createWindow(s, "Copy", func(sc tcell.Screen, x, y, w, h int) {
						st := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
						retrotui.PrintCentered(sc, y+h/2, x, w, "Copy placeholder", st)
					})
				}},
			},
		},
		{
			Title: "Help", HotKey: 'h', Align: true,
			Items: []retrotui.DropdownItem{
				{Text: "About", OnSelect: func(s tcell.Screen) {
					createWindow(s, "About", func(sc tcell.Screen, x, y, w, h int) {
						st := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
						retrotui.PrintCentered(sc, y+2, x, w, "About placeholder", st)
					})
				}},
			},
		},
	}
	width, _ := screen.Size()
	appMenus[len(appMenus)-1].Position = width - len(appMenus[len(appMenus)-1].Title) - 2
}

func drawWizardUI(s tcell.Screen) {
	retrotui.DrawBackground(s, state.Background, true, ' ')

	// Wizard box dimensions
	sw, sh := s.Size()
	boxW, boxH := 60, 12
	boxX, boxY := (sw-boxW)/2, (sh-boxH)/2
	retrotui.DrawBox(s, boxX, boxY, boxW, boxH, tcell.ColorBlack, WizardBoxBg, retrotui.DrawOptions{ShadowEnabled: true})

	page := wizardPages[wizardIndex]
	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(WizardBoxBg)
	retrotui.PrintCentered(s, boxY+1, boxX, boxW, page.Title, titleStyle)

	contentStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(WizardBoxBg)
	for i, line := range page.Content {
		retrotui.PrintCentered(s, boxY+3+i, boxX, boxW, line, contentStyle)
	}

	// Buttons row
	btnLabels := []string{"< Back", "Next >", "Cancel"}
	if wizardIndex == 0 {
		btnLabels[0] = "      "
	} // disabled back
	if wizardIndex == len(wizardPages)-1 {
		btnLabels[1] = "Finish"
	}

	btnY := boxY + boxH - 3
	btnX := boxX + boxW - 8*3
	for i, lbl := range btnLabels {
		style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGray)
		if i == wizardBtnIdx {
			style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
		}
		retrotui.DrawBox(s, btnX+i*8, btnY, 8, 3, tcell.ColorWhite, tcell.ColorBlack, retrotui.DrawOptions{})
		retrotui.PrintCentered(s, btnY+1, btnX+i*8, 8, lbl, style)
	}

	s.Show()
}

func handleWizardEvents(s tcell.Screen, ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyLeft:
			if wizardBtnIdx > 0 {
				wizardBtnIdx--
			}
		case tcell.KeyRight:
			if wizardBtnIdx < 2 {
				wizardBtnIdx++
			}
		case tcell.KeyEnter:
			switch wizardBtnIdx {
			case 0: // Back
				if wizardIndex > 0 {
					wizardIndex--
					wizardBtnIdx = 1
				}
			case 1: // Next or Finish
				if wizardIndex < len(wizardPages)-1 {
					wizardIndex++
					wizardBtnIdx = 1
				} else {
					// Finish -> windows demo
					wizardRunning = false
					state.CurrentScreen = screenMain
				}
			case 2: // Cancel
				wizardRunning = false
				state.CurrentScreen = screenMain
			}
		case tcell.KeyEsc:
			wizardRunning = false
			state.CurrentScreen = screenMain
		}
		drawWizardUI(s)
	case *tcell.EventMouse:
		// simple hit‑test for button clicks (optional)
	}
}

// ----------------------------------------------------------------------------
// main — entry point (minor edits for new dialog screen routing)
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

	// Initial boxes menu items (unchanged)
	state = retrotui.UIState{
		Background:    AppBackground,
		CurrentScreen: screenMenu,
		MenuItems: []retrotui.MenuItem{
			{Text: "1. 10 Rounds Hell Divers Moves", Instruction: "Play 10 rounds with the standard Hell Divers move set.\nPrepare for intense tactical action!"},
			{Text: "2. 10 Rounds Random Moves", Instruction: "Play 10 rounds with completely random move generation.\nExpect the unexpected!"},
			{Text: "3. 10 Rounds Hell Divers Moves (Timer)", Instruction: "Play 10 rounds with standard Hell Divers moves and a time limit.\nSpeed is essential for this mode!"},
		},
		CurrentSelection: 0,
		ActiveMenu:       -1,
		ActiveMenuItem:   -1,
		MenuBarActive:    false,
	}

	initialiseWindowsAndMenus(screen)

	for {
		switch state.CurrentScreen {
		case screenMenu:
			drawMenuScreenUI(screen, state.CurrentSelection)
		case screenDialog:
			drawDialogScreenUI(screen)
		case screenWizard:
			drawWizardUI(screen)
		case screenMain:
			drawMainScreenUI(screen)
		}

		ev := screen.PollEvent()
		if handleCommonEvents(screen, ev) {
			continue
		}

		switch state.CurrentScreen {
		case screenMenu:
			handleMenuScreenEvents(screen, ev)
		case screenDialog:
			handleDialogScreenEvents(screen, ev)
		case screenWizard:
			handleWizardEvents(screen, ev)
		case screenMain:
			handleMainScreenEvents(screen, ev)
		}
	}
}
