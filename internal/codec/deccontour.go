package codec

func DecContour(contour Contour, header *Header) Voxel {
	voxel := make(Voxel, header.Length[0])
	for f := range contour {
		voxel[f] = getFrame(contour[f], header.Length[1], header.Length[2])
	}

	return voxel
}

func getFrame(contourFrame [][]chaincode, height, width int) bitmap {
	frame := make(bitmap, height)
	for y := range frame {
		frame[y] = make([]byte, width)
	}
	segments := make([]bitmap, len(contourFrame))
	for l := range segments {
		segments[l] = make(bitmap, height)
		for y := range segments[l] {
			segments[l][y] = make([]byte, width)
		}
	}

	for l := range segments {
		contourSegment := contourFrame[l]
		if len(contourSegment) > 0 {
			outer := contourSegment[0]
			inners := contourSegment[1:]
			fillSegment(segments[l], outer, inners)
		} else {
			panic("outer contour is not exist")
		}
	}

	for l := range segments {
		for y := range segments[l] {
			for x, v := range segments[l][y] {
				if v == 1 {
					frame[y][x] = 1
				}
			}
		}
	}

	return frame
}

func fillSegment(img bitmap, outer chaincode, inners []chaincode) {
	for _, p := range outer.Points {
		img[p.y][p.x] = 1
	}
	for _, inner := range inners {
		for _, p := range inner.Points {
			img[p.y][p.x] = 2
		}
	}
	for y := range img {
		for x, v := range img[y] {
			if v != 0 {
				// すでに塗り潰した点は除外（線上も除外）
				continue
			}
			p := point{x, y}
			if closedAreaDesicion(p, outer) {
				// もし内輪郭の内側なら 2 で塗り潰し
				isInner := false
				for _, inner := range inners {
					if closedAreaDesicion(p, inner) {
						fillArea(img, p, 0, 2)
						isInner = true
						break
					}
				}
				// 内輪郭の外側で外輪郭の内側ならば 1 になる
				if !isInner {
					img[y][x] = 1
				}
			} else {
				// 外側なら 2 で塗り潰し
				fillArea(img, p, 0, 2)
			}
		}
	}
}
