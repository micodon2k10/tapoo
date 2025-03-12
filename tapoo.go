package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

var count = 0

func drawMaze() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	count++

	if count%2 == 0 {
		termbox.HideCursor()
	} else {
		termbox.SetCursor(30, 30)
	}

	midy := h / 3
	midx := (w - 30) / 3

	midy = 1
	midx = 1

	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+30, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+30, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+30, midy+1, '┘', coldef, coldef)

	for i := 0; i < 30; i++ {
		termbox.SetCell(midx+i, midy-1, '-', coldef, coldef)
	}

	for i := 0; i < 30; i++ {
		termbox.SetCell(midx+i, midy+1, '-', coldef, coldef)
	}

	for i, c := range "Press ESC to quit" {
		termbox.SetCell(midx+6+i+count, midy+3, c, coldef, coldef)
	}

	for i, c := range strconv.Itoa(h) + "= height" {
		termbox.SetCell(midx+6+i, midy+4, c, coldef, coldef)
	}

	for i, c := range strconv.Itoa(w) + "= width" {
		termbox.SetCell(midx+6+i, midy+5, c, coldef, coldef)
	}

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	drawMaze()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
			case termbox.KeyBackspace, termbox.KeyBackspace2:
			case termbox.KeyDelete, termbox.KeyCtrlD:
			case termbox.KeyTab:
			case termbox.KeySpace:
			case termbox.KeyCtrlK:
			case termbox.KeyHome, termbox.KeyCtrlA:
			case termbox.KeyEnd, termbox.KeyCtrlE:
			default:
				if ev.Ch != 0 {
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		drawMaze()

	}
}
