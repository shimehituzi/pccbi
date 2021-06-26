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

func FillArea(img [][]byte, point Point, contour Contour, rect rectMinMax) {
	if img[point.Y][point.X] == 0 && ClosedAreaDesicion(point, contour) {
		fillArea(img, point.X, point.Y, rect)
	}
}

func fillArea(img [][]byte, x, y int, rect rectMinMax) {
	img[y][x] = 2

	if y-1 > rect.minY && img[y-1][x] == 0 {
		fillArea(img, x, y-1, rect)
	}
	if x-1 > rect.minX && img[y][x-1] == 0 {
		fillArea(img, x-1, y, rect)
	}
	if y+1 < rect.maxY && img[y+1][x] == 0 {
		fillArea(img, x, y+1, rect)
	}
	if x+1 < rect.maxX && img[y][x+1] == 0 {
		fillArea(img, x+1, y, rect)
	}
}
