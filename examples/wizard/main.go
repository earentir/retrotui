// examples/wizard/main.go — RetroTUI Wizard Demo
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
	AppBackground = tcell.NewRGBColor(65, 70, 217) // Blue background
	StatusBarFg   = tcell.ColorGreen
	StatusBarBg   = tcell.ColorDarkBlue
	WizardBoxFg   = tcell.ColorWhite
	WizardBoxBg   = tcell.ColorBlue
)

// ----------------------------------------------------------------------------
// Wizard data structure
// ----------------------------------------------------------------------------
type WizardPage struct {
	Title   string
	Content []string
}

var (
	wizardPages = []WizardPage{
		{Title: "Step 1 of 3 — Welcome", Content: []string{"This wizard will guide you", "through basic configuration."}},
		{Title: "Step 2 of 3 — Options", Content: []string{"Choose preferred settings", "and press Next."}},
		{Title: "Step 3 of 3 — Ready", Content: []string{"Setup is complete.", "Click Finish to exit."}},
	}
	wizardIndex  = 0 // current page
	wizardBtnIdx = 1 // 0=Back, 1=Next/Finish, 2=Cancel
)

// ----------------------------------------------------------------------------
// Event handling
// ----------------------------------------------------------------------------
func handleEvents(s tcell.Screen, ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEsc, tcell.KeyCtrlC:
			s.Fini()
			os.Exit(0)
		case tcell.KeyLeft:
			if wizardBtnIdx > 0 {
				wizardBtnIdx--
			}
			drawWizardUI(s)
		case tcell.KeyRight:
			if wizardBtnIdx < 2 {
				wizardBtnIdx++
			}
			drawWizardUI(s)
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
					// At the last page - Finish button
					s.Fini()
					os.Exit(0)
				}
			case 2: // Cancel
				s.Fini()
				os.Exit(0)
			}
			drawWizardUI(s)
		}
	case *tcell.EventMouse:
		mouseX, mouseY := e.Position()
		buttons := e.Buttons()

		if buttons == tcell.ButtonPrimary {
			// Simple button detection
			sw, sh := s.Size()
			boxW, boxH := 60, 12
			boxX, boxY := (sw-boxW)/2, (sh-boxH)/2
			btnY := boxY + boxH - 3
			btnX := boxX + boxW - 8*3

			// Check if clicking on a button
			for i := 0; i < 3; i++ {
				if mouseX >= btnX+i*8 && mouseX < btnX+i*8+8 &&
					mouseY >= btnY && mouseY < btnY+3 {
					wizardBtnIdx = i

					// Simulate pressing Enter on this button
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
							// At the last page - Finish button
							s.Fini()
							os.Exit(0)
						}
					case 2: // Cancel
						s.Fini()
						os.Exit(0)
					}
					drawWizardUI(s)
					break
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
func drawWizardUI(s tcell.Screen) {
	retrotui.DrawBackground(s, AppBackground, true, ' ')

	// Wizard box dimensions
	sw, sh := s.Size()
	boxW, boxH := 60, 12
	boxX, boxY := (sw-boxW)/2, (sh-boxH)/2
	retrotui.DrawBox(s, boxX, boxY, boxW, boxH, tcell.ColorBlack, WizardBoxBg, retrotui.DrawOptions{ShadowEnabled: true})

	// Draw title and content
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
		btnLabels[0] = "      " // Disabled back button
	}
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

	// Status bar
	retrotui.DrawBottomBar(s, "ESC: Cancel", StatusBarFg, StatusBarBg)

	s.Show()
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

	drawWizardUI(screen)

	for {
		ev := screen.PollEvent()
		if !handleEvents(screen, ev) {
			continue
		}
	}
}
