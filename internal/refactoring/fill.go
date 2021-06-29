package refactoring

func getFilledSegments(lbm labeledBitMap, outer []chainCode) ([]labeledBitMap, []rect, []int) {
	segments := make([]labeledBitMap, len(outer))
	rects := make([]rect, len(outer))
	for i := 0; i < len(outer); i++ {
		segments[i] = make(labeledBitMap, len(lbm))
		for y := range segments[i] {
			segments[i][y] = make([]int, len(lbm[y]))
		}
		rects[i] = rect{
			max: point{0, 0},
			min: point{len(lbm[0]), len(lbm)},
		}
	}

	for y := range lbm {
		for x := range lbm[y] {

			label := lbm[y][x] - 1

			if label < 0 {
				continue
			}

			segments[label][y][x] = lbm[y][x]

			if rects[label].max.x < x {
				rects[label].max.x = x
			}
			if rects[label].max.y < y {
				rects[label].max.y = y
			}
			if rects[label].min.x > x {
				rects[label].min.x = x
			}
			if rects[label].min.y > y {
				rects[label].min.y = y
			}
		}
	}

	numOfInners := make([]int, len(outer))
	for label, rect := range rects {
		for y := rect.min.y; y <= rect.max.y; y++ {
			for x := rect.min.x; x <= rect.max.x; x++ {
				if segments[label][y][x] == 0 && closedAreaDesicion(point{x, y}, outer[label]) {
					fillArea(segments[label], rect, x, y, 0, -(label + 1))
					numOfInners[label]++
				}
			}
		}
	}

	for _, segment := range segments {
		for y := range segment {
			for x := range segment[y] {
				if segment[y][x] < lbm[y][x] && lbm[y][x] <= 0 {
					lbm[y][x] = segment[y][x]
				}
			}
		}
	}

	return segments, rects, numOfInners
}

func fillArea(img [][]int, rect rect, x, y, prev, value int) {
	img[y][x] = value

	if y-1 > rect.min.y && img[y-1][x] == prev {
		fillArea(img, rect, x, y-1, prev, value)
	}
	if x-1 > rect.min.x && img[y][x-1] == prev {
		fillArea(img, rect, x-1, y, prev, value)
	}
	if y+1 < rect.max.y && img[y+1][x] == prev {
		fillArea(img, rect, x, y+1, prev, value)
	}
	if x+1 < rect.max.x && img[y][x+1] == prev {
		fillArea(img, rect, x+1, y, prev, value)
	}
}

func closedAreaDesicion(point point, cc chainCode) bool {
	points := cc.points
	length := len(points)
	wn := 0
	for i := 0; i < length-1; i++ {
		if (points[i].y <= point.y) && (points[i+1].y > point.y) {
			vt := (point.y - points[i].y) / (points[i+1].y - points[i].y)
			if point.x < (points[i].x + (vt * (points[i+1].x - points[i].x))) {
				wn++
			}
		} else if (points[i].y > point.y) && (points[i+1].y <= point.y) {
			vt := (point.y - points[i].y) / (points[i+1].y - points[i].y)
			if point.x < (points[i].x + (vt * (points[i+1].x - points[i].x))) {
				wn--
			}
		}
	}
	return wn != 0
}
