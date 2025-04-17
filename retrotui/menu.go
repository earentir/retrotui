package retrotui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// DrawMenuBar draws the menu bar at the top of the screen
func DrawMenuBar(s tcell.Screen, menus []Menu, activeMenu int, menuBarActive bool, menuBarFg, menuBarBg, menuActiveFg, menuActiveBg tcell.Color) {
	width, _ := s.Size()
	menuBarHeight := 1
	menuBarY := 0

	// Fill the menu bar background
	menuOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}
	FillBox(s, 0, menuBarY, width, menuBarHeight, menuBarBg, menuOptions)

	// Draw each menu
	for i, menu := range menus {
		var menuStyle tcell.Style

		// Determine if this menu is the active one
		if menuBarActive && activeMenu == i {
			menuStyle = tcell.StyleDefault.Foreground(menuActiveFg).Background(menuActiveBg)
		} else {
			menuStyle = tcell.StyleDefault.Foreground(menuBarFg).Background(menuBarBg)
		}

		// Determine position based on alignment
		menuX := menu.Position
		if menu.Align { // Right aligned
			menuX = width - len(menu.Title) - 2
		}

		// Draw the menu title with hotkey highlighted
		PrintMenuTitle(s, menuX, menuBarY, menu.Title, menu.HotKey, menuStyle)
	}
}

// DrawDropdownMenu draws a dropdown menu under a menu bar item
func DrawDropdownMenu(s tcell.Screen, menu Menu, activeMenuItem int, dropdownFg, dropdownBg, dropdownActiveFg, dropdownActiveBg, separatorColor tcell.Color) {
	items := menu.Items

	// Find the widest menu item
	maxWidth := 0
	for _, item := range items {
		if len(item.Text) > maxWidth {
			maxWidth = len(item.Text)
		}
	}

	// Add padding
	menuWidth := maxWidth + 4
	menuHeight := len(items) + 2

	// Get position from the menu
	menuX := menu.Position
	menuY := 1 // Just below the menu bar

	// Adjust for right-aligned menus
	if menu.Align {
		menuX = menu.Position - menuWidth + len(menu.Title) + 2
	}

	// Ensure the menu stays within screen bounds
	width, _ := s.Size()
	if menuX+menuWidth > width {
		menuX = width - menuWidth
	}

	// Draw the menu box
	menuOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      true,
		DoubleLine:         false,
	}
	DrawBox(s, menuX, menuY, menuWidth, menuHeight, tcell.ColorBlack, dropdownBg, menuOptions)

	// Draw menu items
	for i, item := range items {
		itemY := menuY + 1 + i
		if item.IsSeparator {
			// Draw a separator line
			for j := 0; j < menuWidth-2; j++ {
				s.SetContent(menuX+1+j, itemY, '─', nil, tcell.StyleDefault.Foreground(separatorColor).Background(dropdownBg))
			}
		} else {
			// Determine style based on whether this item is selected
			itemStyle := tcell.StyleDefault.Foreground(dropdownFg).Background(dropdownBg)
			if activeMenuItem == i {
				itemStyle = tcell.StyleDefault.Foreground(dropdownActiveFg).Background(dropdownActiveBg)
			}

			// Draw the item text
			PrintAt(s, menuX+2, itemY, item.Text, itemStyle)
		}
	}
}

// DrawSelectionDialog draws a central selection dialog with menu items
func DrawSelectionDialog(s tcell.Screen, menuItems []MenuItem, selected int, selectionBg, selectionActiveBg, selectionNumFg, selectionTextFg tcell.Color) {
	width, height := s.Size()

	// Compute dimensions so the dialog fits the menu text
	maxMenuLen := 0
	for _, item := range menuItems {
		if len(item.Text) > maxMenuLen {
			maxMenuLen = len(item.Text)
		}
	}

	// Add padding (2 spaces each side) and top/bottom borders
	dialogWidth := maxMenuLen + 4
	dialogHeight := len(menuItems) + 4

	dialogX := (width - dialogWidth) / 2
	dialogY := (height - dialogHeight) / 2

	// Options for dialog box
	dialogOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      true,
		DoubleLine:         false,
	}

	DrawBox(s, dialogX, dialogY, dialogWidth, dialogHeight, tcell.ColorBlack, selectionBg, dialogOptions)

	// Render menu items inside the dialog
	DrawMenuItems(s, dialogX, dialogY, menuItems, selected, selectionBg, selectionActiveBg, selectionNumFg, selectionTextFg)
}

// DrawMenuItems draws the menu items inside the selection dialog
func DrawMenuItems(s tcell.Screen, dialogX, dialogY int, menuItems []MenuItem, selected int, selectionBg, selectionActiveBg, selectionNumFg, selectionTextFg tcell.Color) {
	// Render menu items inside the dialog with white numbers and black text
	menuStartY := dialogY + 2
	for i, item := range menuItems {
		// Set background color based on selection
		bg := selectionBg
		if i == selected {
			bg = selectionActiveBg // Highlight selected item
		}

		// Find where the number part ends - after the dot
		dotPos := strings.Index(item.Text, ".")

		if dotPos >= 0 && dotPos < len(item.Text)-1 {
			menuItemX := dialogX + 2

			// Just the digit and dot in white
			numPart := item.Text[:dotPos+1]
			PrintAt(s, menuItemX, menuStartY+i, numPart, tcell.StyleDefault.Foreground(selectionNumFg).Background(bg))

			// The space and rest of text in purple
			textPart := item.Text[dotPos+1:]
			PrintAt(s, menuItemX+len(numPart), menuStartY+i, textPart, tcell.StyleDefault.Foreground(selectionTextFg).Background(bg))
		} else {
			// Fallback
			PrintAt(s, dialogX+2, menuStartY+i, item.Text, tcell.StyleDefault.Foreground(selectionTextFg).Background(bg))
		}
	}
}

// DrawInstructionBox draws the instruction box with dynamic text based on selection
func DrawInstructionBox(s tcell.Screen, menuItems []MenuItem, selected int, defaultInstructionText string, instructionFg, instructionBg tcell.Color) {
	width, height := s.Size()

	instrBoxHeight := 4
	instrBoxWidth := 70
	instrBoxX := (width - instrBoxWidth) / 2
	instrBoxY := height - 3 - instrBoxHeight

	// Options for instruction box
	instrOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      true,
		DoubleLine:         false,
	}

	DrawBox(s, instrBoxX, instrBoxY, instrBoxWidth, instrBoxHeight, tcell.ColorWhite, instructionBg, instrOptions)

	// Get the instruction text for the selected menu item
	var instructionText string
	if selected >= 0 && selected < len(menuItems) && menuItems[selected].Instruction != "" {
		instructionText = menuItems[selected].Instruction
	} else {
		instructionText = defaultInstructionText
	}

	// Split instruction text by newlines to handle multi-line instructions
	instructionLines := strings.Split(instructionText, "\n")

	// Calculate starting Y position to center the instruction lines vertically
	startY := instrBoxY + (instrBoxHeight-len(instructionLines))/2

	// Print each line centered in the instruction box
	for i, line := range instructionLines {
		PrintCentered(s, startY+i, instrBoxX, instrBoxWidth, line,
			tcell.StyleDefault.Foreground(instructionFg).Background(instructionBg))
	}
}

// DrawTitleBox draws the title box at the top of the screen
func DrawTitleBox(s tcell.Screen, appName string, copyrightText string, titleFg, titleBg tcell.Color) {
	width, _ := s.Size()
	titleBoxHeight := 4

	// Options for title box
	titleOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      false,
		DoubleLine:         true,
	}

	// Draw a double-line box across the top
	DrawBox(s, 0, 0, width, titleBoxHeight, titleFg, titleBg, titleOptions)
	PrintAt(s, 2, 1, appName, tcell.StyleDefault.Foreground(titleFg).Background(titleBg))
	PrintAt(s, 2, 2, copyrightText, tcell.StyleDefault.Foreground(titleFg).Background(titleBg))
}

// DrawBottomBar draws the status bar at the bottom of the screen
func DrawBottomBar(s tcell.Screen, statusText string, statusFg, statusBg tcell.Color) {
	width, height := s.Size()

	bottomBoxHeight := 1
	bottomBoxY := height - bottomBoxHeight
	bottomOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      false,
	}

	FillBox(s, 0, bottomBoxY, width, bottomBoxHeight, statusBg, bottomOptions)
	PrintAt(s, 2, bottomBoxY, statusText, tcell.StyleDefault.Foreground(statusFg).Background(statusBg))
}

// DrawSimpleMessage draws a simple message box in the center of the screen
func DrawSimpleMessage(s tcell.Screen, message string, fgColor, bgColor tcell.Color) {
	width, height := s.Size()
	msgBoxWidth := len(message) + 6
	msgBoxHeight := 3
	msgBoxX := (width - msgBoxWidth) / 2
	msgBoxY := (height - msgBoxHeight) / 2

	// Options for message box
	msgOptions := DrawOptions{
		FillPatternEnabled: false,
		FillRune:           ' ',
		ShadowEnabled:      true,
		DoubleLine:         false,
	}

	// Draw message box
	DrawBox(s, msgBoxX, msgBoxY, msgBoxWidth, msgBoxHeight, fgColor, bgColor, msgOptions)
	PrintCentered(s, msgBoxY+1, msgBoxX, msgBoxWidth, message, tcell.StyleDefault.Foreground(fgColor).Background(bgColor))

	s.Show()
}

// DrawBackground fills the entire screen with the main background color/pattern
func DrawBackground(s tcell.Screen, bgColor tcell.Color, patternEnabled bool, patternRune rune) {
	options := DrawOptions{
		FillPatternEnabled: patternEnabled,
		FillRune:           patternRune,
		ShadowEnabled:      false,
	}

	bgStyle := tcell.StyleDefault.Background(bgColor).Foreground(tcell.ColorWhite)
	s.Fill(GetFillChar(options), bgStyle)
}
