package labeling

type Point struct {
	X, Y int
}

func NewPoint(X, Y int) Point {
	return Point{X, Y}
}

func (p Point) in(points []Point) bool {
	for _, point := range points {
		if p == point {
			return true
		}
	}
	return false
}

type ChainCode struct {
	Start  Point
	Code   []byte
	Points []Point
}

type Direction struct {
	d    Point
	code byte
	oct  bool
}

func newDirection(code byte, oct bool) Direction {
	direction := [8]Point{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
	var d Point
	if oct {
		d = direction[code]
	} else {
		d = direction[code*2]
	}
	return Direction{d, code, oct}
}

func (d Direction) nextDirections() []Direction {
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
	nextDirections := make([]Direction, numOfDirection)
	for i := range nextDirections {
		nextDirections[i] = newDirection(directionCodes[i], d.oct)
	}
	return nextDirections
}

func ContourTracking(bitmap [][]byte, value byte, oct bool) *ChainCode {
	cc := new(ChainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {

				cc.Start = Point{imageX, imageY}
				cc.Points = []Point{cc.Start}
				currentD := newDirection(0, oct)
				currentP := NewPoint(cc.Start.X, cc.Start.Y)

				checkP := NewCheckPoint(cc.Start.X, cc.Start.Y, bitmap, value, oct)

				for {
					for _, nextD := range currentD.nextDirections() {
						nextP := NewPoint(currentP.X+nextD.d.X, currentP.Y+nextD.d.Y)
						if nextP.Y < 0 || nextP.X < 0 || len(bitmap) <= nextP.Y || len(bitmap[0]) <= nextP.X {
							continue
						}
						if bitmap[nextP.Y][nextP.X] == value {
							cc.Code = append(cc.Code, nextD.code)
							cc.Points = append(cc.Points, nextP)
							currentD = nextD
							currentP = nextP
							break
						}
					}
					if cc.Start == currentP && checkP.in(cc.Points) {
						break
					}
				}

				return cc
			}
		}
	}
	return cc
}

func NewCheckPoint(x, y int, bitmap [][]byte, value byte, oct bool) Point {
	if oct {
		if 0 < x-1 && y+1 < len(bitmap) && bitmap[y+1][x-1] == value {
			return NewPoint(x-1, y+1)
		}
	} else {
		if y+1 < len(bitmap) && bitmap[y+1][x] == value {
			return NewPoint(x, y+1)
		}
	}
	return NewPoint(x, y)
}
