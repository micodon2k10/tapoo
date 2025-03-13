package maze

import (
	"errors"
	"math"
)

// seed defines the size of the maze to be used in the training level (level 0).
// It can also be referred to as the size of the training field.
const seed = 100

// diff defines the difference between maze sizes in consecutive game levels.
const diff = 10

// maxLevel defines the maximum level that can be played in this game.
// Due to the large size of the maze at the final level, it might never be reached especially
// for users with smaller screen sizes.
const maxLevel = 290

// generateMazeArea generates the full maze size depending on the provided game level.
func generateMazeArea(level int) float64 {
	// Level larger than maxLevel should never be used
	if level >= maxLevel {
		level = maxLevel
	}
	return float64((level * diff) + seed)
}

// appendFunc appends the valid dimensions to the size array
func appendFunc(remaider float64, x, y int, tSize Dimensions) []Dimensions {
	items := make([]Dimensions, 0)

	if remaider == 0 && (tSize.Length >= y) && (tSize.Width >= x) {
		items = append(items, Dimensions{Length: y, Width: x})
	}

	if remaider == 0 && (tSize.Length >= x) && (tSize.Width >= y) {
		items = append(items, Dimensions{Length: x, Width: y})
	}

	return items
}

// factorizeMazeArea factorizes the MazeArea using the trial division algorithm
// to get all possible factors for the length and the width values.
// The smallest value of either length or width can only be 5.
func factorizeMazeArea(mazeArea float64, c Dimensions) []Dimensions {
	var size = make([]Dimensions, 0)

	for i := int(math.Sqrt(mazeArea)); i > 4; i-- {
		remaider := math.Remainder(mazeArea, float64(i))
		val := int(mazeArea) / i

		size = append(size, appendFunc(remaider, i, val, c)...)
	}

	return size
}

// getMazeDimensions obtains the best length and width measurements for the
// current level and terminal size provided.
func getMazeDimensions(level int, terminalSize Dimensions) (*Dimensions, error) {
	area := generateMazeArea(level)
	errMsg := "terminal size is too small for the current level"

	if int(area) > (terminalSize.Width * terminalSize.Length) {
		return &Dimensions{}, errors.New(errMsg)
	}

	dimensions := factorizeMazeArea(area, terminalSize)
	totalCount := len(dimensions)

	for i := 0; i < totalCount; i++ {
		return &dimensions[getRandomNo(totalCount)], nil
	}

	// If the terminal size hasn't been minimized, It should never get here
	return &Dimensions{}, errors.New(errMsg)
}

// getTerminalSize calculate the terminal size from the values captured by the
// termbox.Size() function
func getTerminalSize(h, w int) Dimensions {
	return Dimensions{Length: (h - 5) / 4, Width: (w - 10) / 2}
}
