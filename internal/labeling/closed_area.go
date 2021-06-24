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

func FillArea(img [][]byte, point Point, contour Contour) {
	if img[point.Y][point.X] == 0 && ClosedAreaDesicion(point, contour) {
		mx := len(img[0])
		my := len(img)
		fillArea(img, point.X, point.Y, mx, my)
	}
}

func fillArea(img [][]byte, x, y, mx, my int) {
	img[y][x] = 2

	if y-1 > 0 && img[y-1][x] == 1 {
		fillArea(img, x, y-1, mx, my)
	}
	if x-1 > 0 && img[y][x-1] == 1 {
		fillArea(img, x-1, y, mx, my)
	}
	if y+1 < my && img[y+1][x] == 1 {
		fillArea(img, x, y+1, mx, my)
	}
	if x+1 < mx && img[y][x+1] == 1 {
		fillArea(img, x+1, y, mx, my)
	}
}
