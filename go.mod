module dosmenu

go 1.24.1

require github.com/gdamore/tcell/v2 v2.8.1
require github.com/earentir/retrotui v0.0.0-20231009155414-1b2f3a4c5d7e // indirect
replace github.com/earentir/retrotui => ./retrotui

require (
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
