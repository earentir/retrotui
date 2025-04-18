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
