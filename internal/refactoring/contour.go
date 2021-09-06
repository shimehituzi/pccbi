package refactoring

func newContours(orig bitmap) []chainCode {
	return []chainCode{}
}

func closedAreaDesicion(p point, cc chainCode) bool {
	ps := cc.points
	wn := 0
	for i := 0; i < len(ps)-1; i++ {
		if (ps[i].y <= p.y) && (ps[i+1].y > p.y) {
			vt := (p.y - ps[i].y) / (ps[i+1].y - ps[i].y)
			if p.x < (ps[i].x + (vt * (ps[i+1].x - ps[i].x))) {
				wn++
			}
		} else if (ps[i].y > p.y) && (ps[i+1].y <= p.y) {
			vt := (p.y - ps[i].y) / (ps[i+1].y - ps[i].y)
			if p.x < (ps[i].x + (vt * (ps[i+1].x - ps[i].x))) {
				wn--
			}
		}
	}
	return wn != 0
}

func fillArea(img bitmap, p point, prev, value byte) {
	img[p.y][p.x] = value

	nearest4points := [4]point{
		{p.x, p.y - 1},
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
	}

	for _, p := range nearest4points {
		if validPoint(p, img) && img[p.y][p.x] == prev {
			fillArea(img, p, prev, value)
		}
	}
}
