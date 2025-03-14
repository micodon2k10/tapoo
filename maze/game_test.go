package maze

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPlayerMovement tests the functionality of playerMovement
func TestPlayerMovement(t *testing.T) {
	data := [][]string{
		[]string{"|", "---", "|", "---", "|", "---", "|"},
		[]string{"|", " A ", " ", "   ", "|", "   ", "|"},
		[]string{"|", "   ", "|", "   ", "|", "---", "|"},
		[]string{"|", "   ", " ", " B ", " ", "   ", "|"},
		[]string{"|", "---", "|", "   ", "|", "---", "|"},
		[]string{"|", "   ", "|", "   ", "|", "   ", "|"},
		[]string{"|", "---", "|", "---", "|", "---", "|"},
	}

	Convey("TestPlayerMovement: Given the grid view and the current player position", t, func() {
		var d = Dimensions{Length: 3, Width: 3}

		Convey("is at the middle the player should be able to move to all directions"+
			"position exists for the direction provided", func() {
			for direction, output := range map[string][]int{
				"LEFT": []int{3, 1}, "RIGHT": []int{3, 5},
				"DOWN": []int{5, 3}, "UP": []int{1, 3}} {

				d.StartPosition = []int{3, 3}

				d.playerMovement(data, direction)

				So(output[0], ShouldEqual, d.StartPosition[0])
				So(output[1], ShouldEqual, d.StartPosition[1])
			}

			Convey("is at a corner, the player should only be able to move to directions with spaces", func() {
				for direction, output := range map[string][]int{
					"LEFT": []int{1, 1}, "RIGHT": []int{1, 3},
					"DOWN": []int{3, 1}, "UP": []int{1, 1}} {

					d.StartPosition = []int{1, 1}

					d.playerMovement(data, direction)

					So(output[0], ShouldEqual, d.StartPosition[0])
					So(output[1], ShouldEqual, d.StartPosition[1])
				}
			})
		})
	})
}
