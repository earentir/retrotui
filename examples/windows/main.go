// examples/windows/main.go — RetroTUI Windows and Menus Demo
package main

import (
	"os"

	"github.com/earentir/retrotui" // Import from GitHub path
	"github.com/gdamore/tcell/v2"
)

// ----------------------------------------------------------------------------
// Colour scheme
// ----------------------------------------------------------------------------
var (
	AppBackground    = tcell.NewRGBColor(65, 70, 217) // Blue background
	MenuBarFg        = tcell.ColorWhite
	MenuBarBg        = tcell.ColorDarkBlue
	MenuActiveFg     = tcell.ColorWhite
	MenuActiveBg     = tcell.ColorRed
	DropdownFg       = tcell.ColorBlack
	DropdownBg       = tcell.ColorLightGray
	DropdownActiveFg = tcell.ColorWhite
	DropdownActiveBg = tcell.ColorBlue
	SeparatorColor   = tcell.ColorGray
	StatusBarFg      = tcell.ColorGreen
	StatusBarBg      = tcell.ColorDarkBlue
)

// ----------------------------------------------------------------------------
// Application state & globals
// ----------------------------------------------------------------------------
var (
	state      retrotui.UIState
	appMenus   []retrotui.Menu
	appWindows []*retrotui.Window
)

// ----------------------------------------------------------------------------
// Event handling
// ----------------------------------------------------------------------------
func handleEvents(s tcell.Screen, ev tcell.Event) bool {
	// First give top-most window a chance
	if len(appWindows) > 0 {
		if retrotui.ManageWindows(s, appWindows, ev,
			tcell.ColorWhite, tcell.ColorBlue,
			tcell.ColorYellow, tcell.ColorBlue,
			tcell.ColorRed, tcell.ColorBlue) {
			return true
		}
	}

	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyEsc || e.Key() == tcell.KeyF3 || e.Key() == tcell.KeyCtrlC ||
			e.Rune() == 'q' || e.Rune() == 'Q' {
			s.Fini()
			os.Exit(0)
		}

		// Menu keyboard navigation
		if e.Modifiers()&tcell.ModAlt != 0 {
			for i, menu := range appMenus {
				if e.Rune() == menu.HotKey {
					state.ActiveMenu = i
					state.ActiveMenuItem = 0
					state.MenuBarActive = true
					drawUI(s)
					return true
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
			drawUI(s)
		}

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
								return true
							}
						}
					}
					drawUI(s)
					return true
				} else if buttons == tcell.ButtonPrimary {
					state.MenuBarActive = false
					state.ActiveMenu = -1
					drawUI(s)
					return true
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
					drawUI(s)
					return true
				}
			}
			if state.MenuBarActive {
				state.MenuBarActive = false
				state.ActiveMenu = -1
				drawUI(s)
				return true
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
func drawUI(s tcell.Screen) {
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
// Window creation function
// ----------------------------------------------------------------------------
func createWindow(s tcell.Screen, title string, content func(s tcell.Screen, x, y, width, height int)) {
	sw, sh := s.Size()
	w, h := 50, 15
	x, y := (sw-w)/2, (sh-h)/2
	win := retrotui.NewWindow(title, x, y, w, h)
	win.Content = content
	appWindows = append(appWindows, win)
	drawUI(s)
}

// ----------------------------------------------------------------------------
// Menu initialization
// ----------------------------------------------------------------------------
func initialiseMenus(screen tcell.Screen) {
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

	// Initialize application state
	state = retrotui.UIState{
		Background:     AppBackground,
		ActiveMenu:     -1,
		ActiveMenuItem: -1,
		MenuBarActive:  true,
	}

	// Initialize windows and menus
	appWindows = make([]*retrotui.Window, 0)
	initialiseMenus(screen)

	drawUI(screen)

	for {
		ev := screen.PollEvent()
		if !handleEvents(screen, ev) {
			continue
		}
	}
}
