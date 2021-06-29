package refactoring

func getInnerContour(segments []labeledBitMap, rects []rect, numOfInners []int) [][]chainCode {
	innerMatrix := make([][]chainCode, len(segments))
	for label := range innerMatrix {
		innerMatrix[label] = make([]chainCode, numOfInners[label])
		for i := range innerMatrix[label] {
			cc := contourTracking(segments[label], -(label + 1), false)
			innerMatrix[label][i] = cc
			fillArea(segments[label], rects[label], cc.start.x, cc.start.y, -(label + 1), 0)
		}
	}
	return innerMatrix
}
