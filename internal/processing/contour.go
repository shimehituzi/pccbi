package processing

type direction struct {
	d    point
	code byte
	oct  bool
}

func newDirection(code byte, oct bool) direction {
	directions := [8]point{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
	var d point
	if oct {
		d = directions[code]
	} else {
		d = directions[code*2]
	}
	return direction{d, code, oct}
}

func (d direction) nextDirections() []direction {
	var (
		numOfDirection byte
		firstDirection byte
	)
	if d.oct {
		numOfDirection = 8
		firstDirection = d.code + 5
	} else {
		numOfDirection = 4
		firstDirection = d.code + 3
	}
	directionCodes := make([]byte, numOfDirection)
	for i := range directionCodes {
		directionCodes[i] = (firstDirection + byte(i)) % numOfDirection
	}
	nextDirections := make([]direction, numOfDirection)
	for i := range nextDirections {
		nextDirections[i] = newDirection(directionCodes[i], d.oct)
	}
	return nextDirections
}

func contourTracking(bitmap [][]int, value int, oct bool) chainCode {
	cc := *new(chainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {

				cc.start = point{imageX, imageY}
				cc.points = []point{cc.start}
				currentD := newDirection(0, oct)
				currentP := newPoint(cc.start.x, cc.start.y)

				checkP := newCheckPoint(cc.start.x, cc.start.y, bitmap, value, oct)

				for {
					for _, nextD := range currentD.nextDirections() {
						nextP := newPoint(currentP.x+nextD.d.x, currentP.y+nextD.d.y)
						if nextP.y < 0 || nextP.x < 0 || len(bitmap) <= nextP.y || len(bitmap[0]) <= nextP.x {
							continue
						}
						if bitmap[nextP.y][nextP.x] == value {
							cc.code = append(cc.code, nextD.code)
							cc.points = append(cc.points, nextP)
							currentD = nextD
							currentP = nextP
							break
						}
					}
					if cc.start == currentP && checkP.in(cc.points) {
						break
					}
				}

				return cc
			}
		}
	}
	return cc
}

func newCheckPoint(x, y int, bitmap [][]int, value int, oct bool) point {
	if oct {
		if 0 < x-1 && y+1 < len(bitmap) && bitmap[y+1][x-1] == value {
			return newPoint(x-1, y+1)
		}
	} else {
		if y+1 < len(bitmap) && bitmap[y+1][x] == value {
			return newPoint(x, y+1)
		}
	}
	return newPoint(x, y)
}

func newPoint(X, Y int) point {
	return point{X, Y}
}

func (p point) in(points []point) bool {
	for _, point := range points {
		if p == point {
			return true
		}
	}
	return false
}
