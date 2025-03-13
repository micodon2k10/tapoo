package maze

import (
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestMain initializes the test environment
func TestMain(m *testing.M) {
	log.SetFlags(0)

	os.Exit(m.Run())
}

// TestGenerateMaze tests the functionality of GenerateMaze
func TestGenerateMaze(t *testing.T) {
	var (
		compressedView []string

		val = &Dimensions{
			Length: 10,
			Width:  10,
		}
	)

	Convey("Given the correct intensity value", t, func() {
		Convey("The maze should be generated without an error", func() {
			var data, startPos, finalPos, err = val.generateMaze(1)

			So(data, ShouldNotBeEmpty)
			So(startPos, ShouldNotEqual, 0)
			So(finalPos, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			for _, walls := range data {
				compressedView = append(compressedView, strings.Join(walls, ""))
			}

			log.Println("Maze \n", strings.Join(compressedView, ""))
		})
	})
}

// TestCreatePath tests the functionality of createPath
func TestCreatePath(t *testing.T) {
	var (
		val = &Dimensions{
			Length: 7,
			Width:  5,
		}

		gridView, err = val.createPlayingField(1)

		// testData defines the cell number mapped to
		// the coordinates of the common wall.
		testData = map[int][]int{
			10: []int{4, 5},
			16: []int{5, 4},
			18: []int{5, 6},
			24: []int{6, 5},
		}
	)

	Convey("Given a grid view without any common path between cells", t, func() {
		Convey("The common walls should be converted to common paths between cells", func() {
			So(err, ShouldBeNil)
			So(gridView, ShouldNotBeEmpty)

			for cell, loc := range testData {
				val.createPath(gridView[:], 17, cell)

				So(gridView[loc[0]][loc[1]], ShouldContainSubstring, " ")

				log.Println()
				log.Println(gridView)
			}
		})
	})
}

// TestGetPresentNeighbors tests the functionality of getPresentNeighbors
func TestGetPresentNeighbors(t *testing.T) {
	var (
		neighbors []int

		val = &Dimensions{
			Length: 7,
			Width:  5,
		}

		// testData defines the cell number mapped to its unvisited neighbors.
		testData = map[int][]int{
			17: []int{10, 16, 18, 24},
			35: []int{28, 34},
			4:  []int{3, 11, 5},
			7:  []int{6, 14},
		}
	)

	Convey("Given a cell number ", t, func() {
		Convey("The return slice of neighbors should be same as the expected slice", func() {
			for cell, otherCells := range testData {
				neighbors = val.getPresentNeighbors(cell)

				for _, value := range neighbors {
					So(otherCells, ShouldContain, value)
				}
			}
		})
	})
}

// TestGetStartPosition tests the functionality of getStartPosition
func TestGetStartPosition(t *testing.T) {
	var val = &Dimensions{
		Length: 7,
		Width:  5,
	}

	Convey("The start position returned should have less than four neighbors ", t, func() {
		var (
			cellNo    = val.getStartPosition()
			neighbors = val.getPresentNeighbors(cellNo)
		)

		So(len(neighbors), ShouldBeLessThan, 4)

		log.Printf("\nCell : %v \n", cellNo)
		log.Printf("Neighbors : %v \n", neighbors)
	})
}
