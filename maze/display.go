package maze

import (
	"fmt"
	"reflect"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

const (
	intro            = "   You are playing the Maze runner, hide and seek game (Tapoo).      "
	website          = " Visit https://www.tapoo.naihub.com/54ec478gA for more information.  "
	playerNavigation = "      Use the Arrow Keys to navigate the player (in green)           "
	statusMsg        = "         Press Space to Pause.         Scores: %d            "

	space              = "                                                                         "
	pauseMsg           = "                              Game Paused !!!                            "
	gameOverSucceed    = "    Game Over! : Congratulations, Won by Locating the target on time.    "
	gameOverFailed     = "      Game Over! : Ooops!!!, Failed to locate the target on time.        "
	gameOverNavigation = "        Press ESC or Ctrl+C to quit.     Press Ctrl+P to Proceed         "
	highScores         = "                   High Scores: %d                             "
)

// fill prints a string to the termbox view box on the given coordinates.
func fill(x, y int, val string, foreground termbox.Attribute) {
	for index, char := range val {
		termbox.SetCell(x+index, y, char, foreground, coldef)
	}
}

// drawMaze draws the maze on the termbox view.
func drawMaze(config *Dimensions, data [][]string) {
	if err := termbox.Clear(coldef, coldef); err != nil {
		panic(err)
	}

	for loc, msg := range map[int]string{1: intro, 3: website, 5: playerNavigation} {
		fill(len(data[1])/3, loc, msg, coldef)
	}

	for k, d := range data {
		fill(3, 7+k, strings.Join(d, ""), coldef)
	}
}

// refreshUI refreshes the scores value and update the player positions.
func refreshUI(config *Dimensions, count int, data [][]string) {
	drawMaze(config, data)

	termbox.SetCell((targetPos[1]*2)+3, targetPos[0]+7, '#', termbox.ColorRed, termbox.ColorRed)
	termbox.SetCell((startPos[1]*2)+3, startPos[0]+7, '@', termbox.ColorGreen, termbox.ColorGreen)

	fill(len(data[1])/3, len(data)+8, fmt.Sprintf(statusMsg, count), coldef)

	// check if target has been located
	go func() {
		if reflect.DeepEqual(startPos, targetPos) {
			status <- succeeded
		}
	}()

	termbox.Flush()
}

// interruptUI displays some text indicating  if the game is paused or
// after the player won or lost a given tapoo game level.
func interruptUI(msg string, config *Dimensions, data [][]string, color termbox.Attribute) {
	drawMaze(config, data)

	xAxis := len(data[1]) / 4

	for _, loc := range []int{3, 5, 7, 9} {
		fill(xAxis, len(data)/2+loc, space, coldef)
	}

	for loc, msg := range map[int]string{4: msg, 8: gameOverNavigation} {
		fill(xAxis, len(data)/2+loc, msg, coldef)
	}

	scoresMsg := space
	if !paused {
		scoresMsg = fmt.Sprintf(highScores, scores)
	}

	fill(xAxis, len(data)/2+6, scoresMsg, color)

	termbox.Flush()
}
