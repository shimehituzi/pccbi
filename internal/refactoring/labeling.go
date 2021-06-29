package refactoring

// type pointCloudContour struct {
// 	frame []struct {
// 		img      [][]byte
// 		contours []struct {
// 			outer struct {
// 				start  struct{ x, y int }
// 				code   []byte
// 				points []struct{ x, y int }
// 			}
// 			inter []struct {
// 				start  struct{ x, y int }
// 				code   []byte
// 				points []struct{ x, y int }
// 			}
// 			label int
// 		}
// 	}
// 	length [3]int
// }

type labeledPointCloud struct {
	frame  []frame
	length [3]int
}

type frame struct {
	img      [][]byte
	contours []contour
}

type contour struct {
	outer  chainCode
	innter []chainCode
	label  int
}

type chainCode struct {
	start  point
	code   []byte
	points []point
}

type point struct {
	x, y int
}

type labeledBitMap [][]int

func NewLabeledPointCloud(bc *bitCube) *labeledPointCloud {
	pcContour := new(labeledPointCloud)
	pcContour.length = bc.Length
	return pcContour
}

func newLabeledBitMap(img [][]byte) labeledBitMap {

	return make([][]int, 1)
}
