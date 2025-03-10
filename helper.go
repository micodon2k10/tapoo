package tapoo

import (
	"fmt"
	"math"
)

type (
	// CellAddress defines the nine points/coordinates that make up an individual cell.
	// Each of the points define the location of a character that is meant to be a
	// wall or a path of the maze. MiddleCenter represents the part of the path of the maze,
	// while BottomCenter, BottomLeft, BottomRight, MiddleLeft, MiddleRight, TopCenter,
	// TopLeft and TopRight can either be part of the path or part the wall of the maze.
	CellAddress struct {
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

	// CellNeighbors defines the four nieghbors that may surround a given cell.
	// Cells along the maze edges have two to three nieghbors but cells at the center
	// of the maze have four neighbors.
	CellNeighbors struct {
		Bottom int
		Left   int
		Right  int
		Top    int
	}

	// Dimensions defines the actual number of cells that make up the maze along the vertical and
	// the horizontal edges. Length represents the number of the cells along the horizontal
	// edge while Width represents the number of the cells along the vertical edge.
	Dimensions struct {
		Length int
		Width  int
	}
)

// CreatePlayingField creates the initial version of the maze which a grid of cells.
// The cells are created with characters that are printable on the terminal.
// CreatePlayingField accept a paramenter with intensity of thick the maze walls should be
// created.
func (config *Dimensions) CreatePlayingField(intensity int) ([][]string, error) {
	var (
		chars []string
		ok    bool

		data = [][]string{}

		walls = map[int][]string{
			1: []string{"|", "---"},
			2: []string{"╏", "╍╍╍"},
			3: []string{"║", "==="},
		}
	)

	if chars, ok = walls[intensity]; !ok {
		return data, fmt.Errorf(
			"Invalid value of intensity found: %d. Allowed 1, 2 and 3", intensity)
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

// GetCellAddress creates and returns the cell address of the provided cell.
// A cell address is defined by the nine coordinates, where each of them represents the
// actual position of a terminal printable character which becomes part of the maze.
func (config *Dimensions) GetCellAddress(cellNo int) CellAddress {
	var len int

	if cellNo > (config.Length * config.Width) {
		return CellAddress{}
	}

	if len = cellNo % config.Length; len == 0 {
		len = config.Length
	}

	var wid = getCeiledDivisor(cellNo, config.Length) * 2
	len = len * 2

	return CellAddress{
		TopRight:     []int{wid - 2, len - 2},
		TopCenter:    []int{wid - 2, len - 1},
		TopLeft:      []int{wid - 2, len},
		MiddleRight:  []int{wid - 1, len - 2},
		MiddleCenter: []int{wid - 1, len - 1},
		MiddleLeft:   []int{wid - 1, len},
		BottomRight:  []int{wid, len - 2},
		BottomCenter: []int{wid, len - 1},
		BottomLeft:   []int{wid, len},
	}
}

// GetCellNeighbors fetches all the possible neighbors of the provided cell.
func (config *Dimensions) GetCellNeighbors(cellNo int) CellNeighbors {
	if cellNo > (config.Length * config.Width) {
		return CellNeighbors{}
	}

	var (
		right     = cellNo + 1
		left      = cellNo - 1
		top       = cellNo - config.Length
		bottom    = cellNo + config.Length
		neighbors = CellNeighbors{}
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

// getCeiledDivisor calculates the ceiled divisor of the two values passed.
func getCeiledDivisor(num, dinom int) int {
	return int(math.Ceil(float64(num) / float64(dinom)))
}
