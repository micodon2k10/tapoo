package maze

import (
	"fmt"
	"math"
	"time"
)

type (
	// cellAddress defines the nine points/coordinates that make up an individual cell.
	// Each of the points define the location of a character that is meant to be a
	// wall or a path of the maze. MiddleCenter represents a part of the path of the maze,
	// while BottomCenter, BottomLeft, BottomRight, MiddleLeft, MiddleRight, TopCenter,
	// TopLeft and TopRight can either be a part of the path or a part the wall of the maze.
	cellAddress struct {
		BottomCenter []int
		BottomLeft   []int
		BottomRight  []int
		MiddleCenter []int
		MiddleLeft   []int
		MiddleRight  []int
		TopCenter    []int
		TopLeft      []int
		TopRight     []int
	}

	// cellNeighbors defines the four nieghbors that may surround a given cell.
	// Cells along the maze edges have two to three nieghbors but cells at the center
	// of the maze have four neighbors.
	cellNeighbors struct {
		Bottom int
		Left   int
		Right  int
		Top    int
	}
)

// createPlayingField creates the initial version of the maze which is a grid of cells.
// The cells are created with characters that are printable on the terminal.
// createPlayingField accept a paramenter with intensity of how thick the
// maze walls should be.
func (config *Dimensions) createPlayingField(intensity int) ([][]string, error) {
	var (
		chars []string
		err   error

		data = [][]string{}
	)

	if chars, err = getWallCharacters(intensity); err != nil {
		return data, err
	}

	for i := 0; i < (2*config.Width)+1; i++ {
		var val []string

		for k := 0; k < config.Length+1; k++ {
			val = append(val, chars[0])

			switch {
			case k != config.Length && i%2 == 0:
				val = append(val, chars[1])
			case k != config.Length && i%2 != 0:
				val = append(val, "   ")
			default:
				val = append(val, "\n")
			}
		}

		data = append(data, val)
	}
	return data, nil
}

// getCellAddress creates and returns the cell address of the provided cell.
// A cell address is defined by the nine coordinates, where each of them represents the
// actual position of a terminal printable character that becomes a part of the maze.
func (config *Dimensions) getCellAddress(cellNo int) cellAddress {
	var len int

	if cellNo > (config.Length * config.Width) {
		return cellAddress{}
	}

	if len = cellNo % config.Length; len == 0 {
		len = config.Length
	}

	var wid = getCeiledDivisor(cellNo, config.Length) * 2
	len = len * 2

	return cellAddress{
		BottomCenter: []int{wid, len - 1},
		BottomLeft:   []int{wid, len - 2},
		BottomRight:  []int{wid, len},
		MiddleCenter: []int{wid - 1, len - 1},
		MiddleLeft:   []int{wid - 1, len - 2},
		MiddleRight:  []int{wid - 1, len},
		TopCenter:    []int{wid - 2, len - 1},
		TopLeft:      []int{wid - 2, len - 2},
		TopRight:     []int{wid - 2, len},
	}
}

// getCellNeighbors fetches all the possible neighbors of the provided cell.
func (config *Dimensions) getCellNeighbors(cellNo int) cellNeighbors {
	if cellNo > (config.Length * config.Width) {
		return cellNeighbors{}
	}

	var (
		right     = cellNo + 1
		left      = cellNo - 1
		top       = cellNo - config.Length
		bottom    = cellNo + config.Length
		neighbors = cellNeighbors{}
	)

	if getCeiledDivisor(right, config.Length) == getCeiledDivisor(cellNo, config.Length) {
		neighbors.Right = right
	}

	if getCeiledDivisor(left, config.Length) == getCeiledDivisor(cellNo, config.Length) {
		neighbors.Left = left
	}

	if top > 0 {
		neighbors.Top = top
	}

	if bottom <= (config.Length * config.Width) {
		neighbors.Bottom = bottom
	}

	return neighbors
}

// getRandomNo returns a random number generated from
// the current timestamp and should be less the max value
// provided and greater than or equal to zero. (0 <= X < max)
func getRandomNo(max int) int {
	return int(time.Now().UnixNano() % int64(max))
}

// getCeiledDivisor calculates the ceiled divisor of the two values passed.
func getCeiledDivisor(num, dinom int) int {
	return int(math.Ceil(float64(num) / float64(dinom)))
}

// getWallCharacters returns the maze wall characters associated with the provided intensity.
// If invalid intensity is used an error is thrown.
func getWallCharacters(intensity int) ([]string, error) {
	chars, ok := map[int][]string{
		1: []string{"|", "---", "-"},
		2: []string{"╏", "╍╍╍", "╍"},
		3: []string{"║", "===", "="},
	}[intensity]

	if ok {
		return chars, nil
	}

	return chars, fmt.Errorf(
		"Invalid value of intensity found: %d. Allowed 1, 2 and 3", intensity)
}
