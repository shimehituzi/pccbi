package refactoring

type rect struct {
	max point
	min point
}

func fillLabeledBitMap(lbm labeledBitMap, outer []chainCode) {
	numOfOuter := len(outer)

	imgs := make([]labeledBitMap, numOfOuter)
	rects := make([]rect, numOfOuter)
	for i := 0; i < numOfOuter; i++ {
		imgs[i] = make(labeledBitMap, len(lbm))
		for y := range imgs[i] {
			imgs[i][y] = make([]int, len(lbm[y]))
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

			imgs[label][y][x] = lbm[y][x]

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

	numOfInner := 0
	for label, rect := range rects {
		for y := rect.min.y; y <= rect.max.y; y++ {
			for x := rect.min.x; x <= rect.max.x; x++ {
				if fillArea(imgs[label], point{x, y}, outer[label], rect, 0, -(label + 1)) {
					numOfInner++
				}
			}
		}
	}

	for _, img := range imgs {
		for y := range img {
			for x := range img[y] {
				if img[y][x] < lbm[y][x] && lbm[y][x] <= 0 {
					lbm[y][x] = img[y][x]
				}
			}
		}
	}
}

func fillArea(img [][]int, point point, cc chainCode, rect rect, prev, value int) bool {
	if img[point.y][point.x] == 0 && closedAreaDesicion(point, cc) {
		_fillArea(img, rect, point.x, point.y, prev, value)
		return true
	}
	return false
}

func _fillArea(img [][]int, rect rect, x, y, prev, value int) {
	img[y][x] = value

	if y-1 > rect.min.y && img[y-1][x] == prev {
		_fillArea(img, rect, x, y-1, prev, value)
	}
	if x-1 > rect.min.x && img[y][x-1] == prev {
		_fillArea(img, rect, x-1, y, prev, value)
	}
	if y+1 < rect.max.y && img[y+1][x] == prev {
		_fillArea(img, rect, x, y+1, prev, value)
	}
	if x+1 < rect.max.x && img[y][x+1] == prev {
		_fillArea(img, rect, x+1, y, prev, value)
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
