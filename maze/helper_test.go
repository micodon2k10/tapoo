package maze

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestCreatePlayingField tests the functionality of CreatePlayingField
func TestCreatePlayingField(t *testing.T) {
	var (
		compressedView []string
		err            error
		gridView       [][]string

		val = &Dimensions{
			Length: 17,
			Width:  5,
		}
	)

	Convey("Given the the maze wall intensity ", t, func() {
		Convey("A grid view should be generated without an error given "+
			"the correct intensity ", func() {
			for _, intensity := range []int{1, 2, 3} {
				gridView, err = val.CreatePlayingField(intensity)

				So(err, ShouldBeNil)
				So(gridView, ShouldNotBeEmpty)

				for _, line := range gridView {
					So(line, ShouldHaveLength, (val.Length+1)*2)

					compressedView = append(compressedView, strings.Join(line, ""))
				}

				log.Println("\n", strings.Join(compressedView, ""))

				compressedView = []string{}
			}

		})

		Convey("An error should be thrown given invalid intensity", func() {
			gridView, err = val.CreatePlayingField(10)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Invalid value of intensity found: 10")
			So(gridView, ShouldBeEmpty)
		})
	})
}

// TestGetCellAddress tests the functionality of GetCellAddress
func TestGetCellAddress(t *testing.T) {
	var (
		k CellAddress

		// data represents an initial version of the maze of 6 cells by 5 cells with
		// the respective cells labeled with their cell numbers.
		data = [][]string{
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
			[]string{"|", " 1 ", "|", " 2 ", "|", " 3 ", "|", " 4 ", "|", " 5 ", "|", " 6 ", "|", "\n"},
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
			[]string{"|", " 7 ", "|", " 8 ", "|", " 9 ", "|", "10 ", "|", "11 ", "|", "12 ", "|", "\n"},
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
			[]string{"|", "13 ", "|", "14 ", "|", "15 ", "|", "16 ", "|", "17 ", "|", "18 ", "|", "\n"},
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
			[]string{"|", "19 ", "|", "20 ", "|", "21 ", "|", "22 ", "|", "23 ", "|", "24 ", "|", "\n"},
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
			[]string{"|", "25 ", "|", "26 ", "|", "27 ", "|", "28 ", "|", "29 ", "|", "30 ", "|", "\n"},
			[]string{"|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "---", "|", "\n"},
		}

		val = &Dimensions{
			Length: 6,
			Width:  5,
		}
	)

	Convey("Given the grid view of 6 by 5 cells of a maze ", t, func() {

		Convey("The cell address returned should help fetch a cell label"+
			" matching the cell number ", func() {
			for i := 1; i <= (val.Length * val.Width); i++ {
				k = val.GetCellAddress(i)

				So(data[k.MiddleCenter[0]][k.MiddleCenter[1]], ShouldContainSubstring, strconv.Itoa(i))

				fmt.Println("\ncell", i)
				fmt.Println(data[k.TopRight[0]][k.TopRight[1]], data[k.TopCenter[0]][k.TopCenter[1]], data[k.TopLeft[0]][k.TopLeft[1]])
				fmt.Println(data[k.MiddleRight[0]][k.MiddleRight[1]], data[k.MiddleCenter[0]][k.MiddleCenter[1]], data[k.MiddleLeft[0]][k.MiddleLeft[1]])
				fmt.Println(data[k.BottomRight[0]][k.BottomRight[1]], data[k.BottomCenter[0]][k.BottomCenter[1]], data[k.BottomLeft[0]][k.BottomLeft[1]])
			}
		})

	})
}

// TestGetCellNeighbors tests the functionality of GetCellNeighbors
func TestGetCellNeighbors(t *testing.T) {
	var (
		key             int
		found, expected CellNeighbors

		// data defines the respective neighbors of the provided cell.
		// A neighbor can be on Top, Left, Bottom or Right of the provided cell.
		data = map[int]CellNeighbors{
			1:  CellNeighbors{Top: 0, Right: 2, Bottom: 7, Left: 0},
			2:  CellNeighbors{Top: 0, Right: 3, Bottom: 8, Left: 1},
			3:  CellNeighbors{Top: 0, Right: 4, Bottom: 9, Left: 2},
			4:  CellNeighbors{Top: 0, Right: 5, Bottom: 10, Left: 3},
			5:  CellNeighbors{Top: 0, Right: 6, Bottom: 11, Left: 4},
			6:  CellNeighbors{Top: 0, Right: 0, Bottom: 12, Left: 5},
			7:  CellNeighbors{Top: 1, Right: 8, Bottom: 13, Left: 0},
			8:  CellNeighbors{Top: 2, Right: 9, Bottom: 14, Left: 7},
			9:  CellNeighbors{Top: 3, Right: 10, Bottom: 15, Left: 8},
			10: CellNeighbors{Top: 4, Right: 11, Bottom: 16, Left: 9},
			11: CellNeighbors{Top: 5, Right: 12, Bottom: 17, Left: 10},
			12: CellNeighbors{Top: 6, Right: 0, Bottom: 18, Left: 11},
			13: CellNeighbors{Top: 7, Right: 14, Bottom: 19, Left: 0},
			14: CellNeighbors{Top: 8, Right: 15, Bottom: 20, Left: 13},
			15: CellNeighbors{Top: 9, Right: 16, Bottom: 21, Left: 14},
			16: CellNeighbors{Top: 10, Right: 17, Bottom: 22, Left: 15},
			17: CellNeighbors{Top: 11, Right: 18, Bottom: 23, Left: 16},
			18: CellNeighbors{Top: 12, Right: 0, Bottom: 24, Left: 17},
			19: CellNeighbors{Top: 13, Right: 20, Bottom: 25, Left: 0},
			20: CellNeighbors{Top: 14, Right: 21, Bottom: 26, Left: 19},
			21: CellNeighbors{Top: 15, Right: 22, Bottom: 27, Left: 20},
			22: CellNeighbors{Top: 16, Right: 23, Bottom: 28, Left: 21},
			23: CellNeighbors{Top: 17, Right: 24, Bottom: 29, Left: 22},
			24: CellNeighbors{Top: 18, Right: 0, Bottom: 30, Left: 23},
			25: CellNeighbors{Top: 19, Right: 26, Bottom: 0, Left: 0},
			26: CellNeighbors{Top: 20, Right: 27, Bottom: 0, Left: 25},
			27: CellNeighbors{Top: 21, Right: 28, Bottom: 0, Left: 26},
			28: CellNeighbors{Top: 22, Right: 29, Bottom: 0, Left: 27},
			29: CellNeighbors{Top: 23, Right: 30, Bottom: 0, Left: 28},
			30: CellNeighbors{Top: 24, Right: 0, Bottom: 0, Left: 29},
			31: CellNeighbors{Top: 0, Right: 0, Bottom: 0, Left: 0},
		}

		val = &Dimensions{
			Length: 6,
			Width:  5,
		}
	)

	Convey("Given cells with their expected neighbors in a grid view of 6 by 5 cells", t, func() {

		Convey("The fetched neighbors of the provided cell number "+
			"should match the expected neighbors", func() {

			for key, expected = range data {
				found = val.GetCellNeighbors(key)

				So(expected.Top, ShouldEqual, found.Top)
				So(expected.Left, ShouldEqual, found.Left)
				So(expected.Bottom, ShouldEqual, found.Bottom)
				So(expected.Right, ShouldEqual, found.Right)
			}
		})
	})
}
