package codec

func newChainCode(img bitmap, value byte, inner bool) *chainCode {
	for y := range img {
		for x, v := range img[y] {
			if v == value {
				start := point{x, y}
				return contourTracking(img, start, value, inner)
			}
		}
	}
	return nil
}

func contourTracking(img bitmap, start point, value byte, inner bool) *chainCode {
	cc := new(chainCode)
	cc.start = start
	cc.points = []point{start}

	currentD := newDirection(0)
	currentP := point{start.x, start.y}

	checkP := point{start.x - 1, start.y + 1}
	if !(validPoint(checkP, img) && img[checkP.y][checkP.x] == value) {
		checkP = start
	}

	divisor := byte(8)

	for {
		for _, nextD := range currentD.nextDirections() {
			nextP := point{currentP.x + nextD.d.x, currentP.y + nextD.d.y}
			if validPoint(nextP, img) && img[nextP.y][nextP.x] == value {
				if inner && (nextD.code%2) == 1 {
					beforeD := newDirection((nextD.code - 1) % 8)
					beforeP := point{currentP.x + beforeD.d.x, currentP.y + beforeD.d.y}
					afterD := newDirection((nextD.code + 1) % 8)
					afterP := point{currentP.x + afterD.d.x, currentP.y + afterD.d.y}
					if img[beforeP.y][beforeP.x] == 1 && img[afterP.y][afterP.x] == 1 {
						continue
					}
				}
				cc.code = append(cc.code, (nextD.code-currentD.code)%divisor)
				cc.points = append(cc.points, nextP)
				currentD = nextD
				currentP = nextP
				break
			}
		}
		if start == currentP && checkP.in(cc.points) {
			break
		}
	}

	return cc
}

func newDirection(code byte) direction {
	directions := [8]point{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
	d := directions[code]
	return direction{d, code}
}

func (d direction) nextDirections() []direction {
	numOfDirection := byte(8)
	firstDirection := byte(d.code + 5)
	directionCodes := make([]byte, numOfDirection)
	for i := range directionCodes {
		directionCodes[i] = (firstDirection + byte(i)) % numOfDirection
	}
	nextDirections := make([]direction, numOfDirection)
	for i := range nextDirections {
		nextDirections[i] = newDirection(directionCodes[i])
	}
	return nextDirections
}

func validPoint(p point, img bitmap) bool {
	if p.y < 0 || p.x < 0 || len(img) <= p.y || len(img[0]) <= p.x {
		return false
	}
	return true
}

func (p point) in(points []point) bool {
	for _, point := range points {
		if p == point {
			return true
		}
	}
	return false
}
