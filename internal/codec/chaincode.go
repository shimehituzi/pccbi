package codec

func newChaincode(img bitmap, start point, value byte, inner bool) *chaincode {
	cc := new(chaincode)
	cc.start = start
	cc.points = []point{start}

	currentD := newDirection(0)
	currentP := point{start.x, start.y}

	checkP := point{start.x - 1, start.y + 1}
	if !(validPointByte(checkP, img) && img[checkP.y][checkP.x] == value) {
		checkP = start
	}

	divisor := byte(8)

	for {
		for _, nextD := range currentD.nextDirections() {
			nextP := point{currentP.x + nextD.d.x, currentP.y + nextD.d.y}
			if validPointByte(nextP, img) && img[nextP.y][nextP.x] == value {
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
