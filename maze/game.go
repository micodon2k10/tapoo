package maze

import (
	"fmt"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// coldef maintains the original color used on the
// background or the foreground depending on its usage.
const coldef = termbox.ColorDefault

const (
	// succeeded status is update ONLY when the player locates the target successfully.
	succeeded = iota

	// proceed status should be updated if the player wants to continue playing the game after:
	// 1. They paused the game.
	// 2. They successfully located the target on time and would like to play the next level.
	// 3. They failed to locate on time and would like to play the level again
	proceed

	// failed status is updated after the player fails to locate the target of time.
	failed

	// pause status is updated after the play voluntarity stops the game.
	pause

	// Status should be updated after the player voluntarily paused the game or won the level
	// or even failed to finish the level successfully.
	quit
)

var (
	scores int

	paused = false

	status = make(chan int)
)

// playerMovement calculates the actual player position
// depending on the navigation keys pressed.
func (config *Dimensions) playerMovement(data [][]string, direction string) {
	startPos := config.StartPosition
	xVal, zVal := startPos[1], startPos[0]

	switch {
	case (direction == "LEFT") && ((xVal - 2) > 0) && isSpaceFound(data[zVal][xVal-1]):
		startPos[1] = xVal - 2

	case (direction == "RIGHT") && ((xVal + 2) <= config.Length*2) && isSpaceFound(data[zVal][xVal+1]):
		startPos[1] = xVal + 2

	case (direction == "UP") && ((zVal - 2) > 0) && isSpaceFound(data[zVal-1][xVal]):
		startPos[0] = zVal - 2

	case (direction == "DOWN") && ((zVal + 2) <= config.Width*2) && isSpaceFound(data[zVal+1][xVal]):
		startPos[0] = zVal + 2
	}
}

// handlePlayerMovement detects that keys pressed on the keyboard
// and provides that direction that the player should move to.
func (config *Dimensions) handlePlayerMovement(event termbox.Key, data [][]string) {
	switch event {
	case termbox.KeyEsc, termbox.KeyCtrlC:
		status <- quit

	case termbox.KeyCtrlP:
		status <- proceed

	case termbox.KeySpace:
		status <- pause

	case termbox.KeyArrowLeft:
		config.playerMovement(data, "LEFT")

	case termbox.KeyArrowRight:
		config.playerMovement(data, "RIGHT")

	case termbox.KeyArrowUp:
		config.playerMovement(data, "UP")

	case termbox.KeyArrowDown:
		config.playerMovement(data, "DOWN")
	}
}

// handleKeyboardMapping handles all the keyboard input as captured by termbox
func (config *Dimensions) handleKeyboardMapping(data [][]string) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			config.handlePlayerMovement(ev.Key, data)

		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// Start define where the tapoo game starts at.
func Start() {
	var (
		data [][]string

		// If an error is thrown print it and exit
		errfunc = func(err error) {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	)

	errfunc(termbox.Init())

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	val, err := getMazeDimensions(1, getTerminalSize(termbox.Size()))
	errfunc(err)

	data, err = val.generateMaze(1)
	errfunc(err)

	go val.handleKeyboardMapping(data)

	var (
		timer       = time.NewTicker(500 * time.Microsecond)
		totalCells  = val.Length * val.Width
		timeout     = time.NewTimer(time.Duration(totalCells) * time.Second)
		currentTime = time.Now().Unix()
	)

mainloop:
	for {
		select {
		case timeVal := <-timer.C:
			scores = (totalCells - int(timeVal.Unix()-currentTime)) * 100

			refreshUI(val, scores, data)
		case <-timeout.C:
			go func() { status <- failed }()

		case returnedStatus := <-status:
			timer.Stop()
			timeout.Stop()

			switch {
			case returnedStatus == succeeded:
				interruptUI(gameOverSucceed, val, data, termbox.ColorGreen)
				paused = true

			case returnedStatus == failed:
				interruptUI(gameOverFailed, val, data, termbox.ColorRed)
				paused = true

			case returnedStatus == quit && paused:
				break mainloop

			case returnedStatus == proceed && paused:
				// paused = false

			case returnedStatus == pause:
				paused = true
				interruptUI(pauseMsg, val, data, termbox.ColorBlue)
			}
		}
	}
}
