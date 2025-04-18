package retrotui

import (
	"github.com/gdamore/tcell/v2"
)

// DrawOptions struct for UI components
type DrawOptions struct {
	FillPatternEnabled bool        // true for fill pattern, false for no fill
	FillRune           rune        // fill pattern character ░ or ▒ or ▓ or space
	ShadowEnabled      bool        // true for shadow, false for no shadow
	ShadowRune         rune        // shadow character █
	ShadowColor        tcell.Color // color of the shadow
	DoubleLine         bool        // true for double line box, false for single line
}

// MenuItem struct to hold menu item text and its instruction
type MenuItem struct {
	Text        string
	Instruction string
}

// Menu represents a drop-down menu
type Menu struct {
	Title    string
	HotKey   rune // The key after Alt to activate this menu
	Items    []DropdownItem
	Position int  // x position on menu bar
	Align    bool // false = left, true = right
}

// DropdownItem represents an item within a dropdown menu
type DropdownItem struct {
	Text        string
	IsSeparator bool
	OnSelect    func(s tcell.Screen)
}

// UIState holds the current state of the UI
type UIState struct {
	Background       tcell.Color
	ActiveMenu       int // -1 for no active menu
	ActiveMenuItem   int
	CurrentScreen    string // Identifies which screen is currently active
	MenuBarActive    bool
	MenuItems        []MenuItem // For the menu screen
	CurrentSelection int        // For the menu screen
}

// MenuConfig contains all configuration options for a selection menu
type MenuConfig struct {
	// Application information
	AppName       string
	CopyrightText string

	// Menu items
	MenuItems []MenuItem

	// Default instruction text
	DefaultInstructionText string

	// Color scheme
	Background        tcell.Color
	TitleBarFg        tcell.Color
	TitleBarBg        tcell.Color
	StatusBarFg       tcell.Color
	StatusBarBg       tcell.Color
	MainSelectionBg   tcell.Color
	SelectionActiveBg tcell.Color
	SelectionNumFg    tcell.Color
	SelectionTextFg   tcell.Color
	InstructionBoxFg  tcell.Color
	InstructionBoxBg  tcell.Color

	// Key bindings
	ExitKeys  []tcell.Key
	ExitRunes []rune

	// Behavior options
	ReturnToMenuAfterSelection bool
}

// NavigationAction represents a navigation action
type NavigationAction int

const (
	NavNone NavigationAction = iota
	NavUp
	NavDown
	NavLeft
	NavRight
	NavSelect
	NavExit
)

// KeyConfig holds configuration for keyboard shortcuts
type KeyConfig struct {
	ExitKeys    []tcell.Key // Keys that will exit the application
	ExitRunes   []rune      // Runes that will exit the application
	NavUpKey    tcell.Key   // Key for navigating up
	NavDownKey  tcell.Key   // Key for navigating down
	NavLeftKey  tcell.Key   // Key for navigating left
	NavRightKey tcell.Key   // Key for navigating right
	SelectKey   tcell.Key   // Key for selection
}
