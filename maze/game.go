package maze

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

var (
	scores              int
	startPos, targetPos []int

	exit = make(chan int)
)

func fill(x, y int, val string) {
	for index, char := range val {
		termbox.SetCell(x+index, y, char, coldef, coldef)
	}
}

// drawMaze draws the maze
func drawMaze(config *Dimensions, data [][]string) {
	var err = termbox.Clear(coldef, coldef)
	if err != nil {
		panic(err)
	}

	fill(len(data[1])/3, 1, "You are playing the Maze runner, hide and seek game (Tapoo).")
	fill(len(data[1])/2, 3, "Visit www.tapoo.com for more information.")
	fill(len(data[1])/2, 5, "Use the Arrow Keys to navigate the player (in green)")

	for k, d := range data {
		fill(3, 7+k, strings.Join(d, ""))
	}
}

// refreshUI refreshes the score value and the blinking cursor
func refreshUI(config *Dimensions, count int, data [][]string) {
	drawMaze(config, data)

	termbox.SetCell((targetPos[1]*2)+3, targetPos[0]+7, '#', termbox.ColorRed, termbox.ColorRed)
	termbox.SetCell((startPos[1]*2)+3, startPos[0]+7, '@', termbox.ColorGreen, termbox.ColorGreen)

	fill(len(data[1])/2, len(data)+8, "Press ESC or Ctrl+C to quit.         Scores: "+strconv.Itoa(count))

	// check if target has been located
	go func() {
		if reflect.DeepEqual(startPos, targetPos) {
			exit <- 0
		}
	}()

	termbox.Flush()
}

func gameOverUI(msg string, config *Dimensions, data [][]string) {
	drawMaze(config, data)

	fill(len(data[1])/3, len(data)/2+3, "                                                         ")
	fill(len(data[1])/3, len(data)/2+4, "    Game Over! : "+msg)
	fill(len(data[1])/3, len(data)/2+5, "                                                         ")
	fill(len(data[1])/3, len(data)/2+6, "              High Scores: "+strconv.Itoa(scores)+"                        ")
	fill(len(data[1])/3, len(data)/2+7, "                                                         ")
	fill(len(data[1])/3, len(data)/2+8, "     Press ESC or Ctrl+C to quit.                             ")
	fill(len(data[1])/3, len(data)/2+9, "                                                           ")

	termbox.Flush()
}

func handlePlayersMovement(config *Dimensions, data [][]string) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			var (
				xVal = startPos[1]
				zVal = startPos[0]
			)
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				exit <- 2

			case termbox.KeyArrowLeft:
				if (xVal-2) > 0 && strings.Contains(data[zVal][xVal-1], " ") {
					startPos[1] = xVal - 2
				}

			case termbox.KeyArrowRight:
				if (xVal+2) <= config.Length*2 && strings.Contains(data[zVal][xVal+1], " ") {
					startPos[1] = xVal + 2
				}

			case termbox.KeyArrowUp:
				if (zVal-2) > 0 && strings.Contains(data[zVal-1][xVal], " ") {
					startPos[0] = zVal - 2
				}

			case termbox.KeyArrowDown:
				if (zVal+2) <= config.Width*2 && strings.Contains(data[zVal+1][xVal], " ") {
					startPos[0] = zVal + 2
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// Start define where the tapoo game
// execution should start from.
func Start() {
	var (
		data                  [][]string
		startCell, targetCell int

		timer      = time.NewTicker(500 * time.Microsecond)
		val        = &Dimensions{Length: 30, Width: 7}
		totalCells = val.Length * val.Width
		timeout    = time.NewTimer(time.Duration(totalCells) * time.Second)

		currentTime = time.Now().Unix()
		err         = termbox.Init()
	)

	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	data, startCell, targetCell, err = val.GenerateMaze(1)
	if err != nil {
		panic(err)
	}

	startPos = val.GetCellAddress(startCell).MiddleCenter
	targetPos = val.GetCellAddress(targetCell).MiddleCenter

	go handlePlayersMovement(val, data)

mainloop:
	for {
		select {
		case timeVal := <-timer.C:
			scores = (totalCells - int(timeVal.Unix()-currentTime)) * 100

			refreshUI(val, scores, data)
		case <-timeout.C:
			go func() {
				exit <- 1
			}()

		case exitStatus := <-exit:
			timer.Stop()
			timeout.Stop()
			switch exitStatus {
			case 0:
				gameOverUI("You won after locating the target on time. ", val, data)
			case 1:
				gameOverUI("You failed to locate the target on time. ", val, data)
			case 2:
				break mainloop
			}
		}
	}
}
