package labeling

func ClosedAreaDesicion(point Point, contour Contour) bool {
	points := contour.ChainCode.Points
	length := len(points)
	wn := 0
	for i := 0; i < length-1; i++ {
		if (points[i].Y <= point.Y) && (points[i+1].Y > point.Y) {
			vt := (point.Y - points[i].Y) / (points[i+1].Y - points[i].Y)
			if point.X < (points[i].X + (vt * (points[i+1].X - points[i].X))) {
				wn++
			}
		} else if (points[i].Y > point.Y) && (points[i+1].Y <= point.Y) {
			vt := (point.Y - points[i].Y) / (points[i+1].Y - points[i].Y)
			if point.X < (points[i].X + (vt * (points[i+1].X - points[i].X))) {
				wn--
			}
		}
	}
	return wn != 0
}

func FillArea(img [][]byte, point Point, contour Contour, rect rectMinMax) bool {
	if img[point.Y][point.X] == 0 && ClosedAreaDesicion(point, contour) {
		fillArea(img, point.X, point.Y, rect, 0, 2)
		return true
	}
	return false
}

func fillArea(img [][]byte, x, y int, rect rectMinMax, prev, value byte) {
	img[y][x] = value

	if y-1 > rect.minY && img[y-1][x] == prev {
		fillArea(img, x, y-1, rect, prev, value)
	}
	if x-1 > rect.minX && img[y][x-1] == prev {
		fillArea(img, x-1, y, rect, prev, value)
	}
	if y+1 < rect.maxY && img[y+1][x] == prev {
		fillArea(img, x, y+1, rect, prev, value)
	}
	if x+1 < rect.maxX && img[y][x+1] == prev {
		fillArea(img, x+1, y, rect, prev, value)
	}
}
