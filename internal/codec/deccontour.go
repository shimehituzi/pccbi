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
	outerPoints := outer.getPoints()
	for _, p := range outerPoints {
		img[p.y][p.x] = 1
	}
	for _, p := range outerPoints {
		nearest4points := [4]point{
			{p.x, p.y - 1},
			{p.x - 1, p.y},
			{p.x + 1, p.y},
			{p.x, p.y + 1},
		}
		for _, p := range nearest4points {
			if validPointByte(p, img) && img[p.y][p.x] == 0 && closedAreaDesicion(p, outerPoints) {
				fillArea(img, p, 0, 1)
			}
		}
	}
	innerPointsArray := make([][]point, len(inners))
	for i, inner := range inners {
		innerPointsArray[i] = inner.getPoints()
	}
	for _, innerPoints := range innerPointsArray {
		for _, p := range innerPoints {
			img[p.y][p.x] = 2
		}
	}
	for _, innerPoints := range innerPointsArray {
		for _, p := range innerPoints {
			nearest4points := [4]point{
				{p.x, p.y - 1},
				{p.x - 1, p.y},
				{p.x + 1, p.y},
				{p.x, p.y + 1},
			}
			for _, p := range nearest4points {
				if validPointByte(p, img) && img[p.y][p.x] == 1 && closedAreaDesicion(p, innerPoints) {
					fillArea(img, p, 1, 2)
				}
			}
		}
	}
}
