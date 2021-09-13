package encoder

func newContour(orig bitmap) contour {
	img := make(bitmap, len(orig))
	for i := range orig {
		img[i] = make([]byte, len(orig[i]))
		copy(img[i], orig[i])
	}

	cont := contour{}

	// 外輪郭
	outer := newChainCode(img, 1, false)
	if outer == nil {
		panic("cannot produce chainCode")
	}
	cont = append(cont, *outer)

	// 塗り潰し
	filledOutside := false
	label := byte(2)

	for y := range img {
		for x, v := range img[y] {
			if v == 0 {
				p := point{x, y}
				if filledOutside {
					fillArea(img, p, 0, label)
					label++
				} else {
					if closedAreaDesicion(p, *outer) {
						fillArea(img, p, 0, label)
						label++
					} else {
						fillArea(img, p, 0, 1)
						filledOutside = true
					}
				}
			}
		}
	}

	// 内輪郭
	for l := byte(2); l < label; l++ {
		inner := newChainCode(img, l, true)
		if outer == nil {
			panic("cannot produce chainCode")
		}
		cont = append(cont, *inner)
	}

	return cont
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
