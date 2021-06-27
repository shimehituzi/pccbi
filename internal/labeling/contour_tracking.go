package labeling

type Point struct {
	X, Y int
}

func NewPoint(X, Y int) *Point {
	return &Point{X, Y}
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
		firstDirection = 5
	} else {
		numOfDirection = 4
		firstDirection = 3
	}
	directionCodes := make([]byte, numOfDirection)
	for i := range directionCodes {
		directionCodes[i] = (firstDirection + d.code + byte(i)) % numOfDirection
	}
	nextDirections := make([]Direction, numOfDirection)
	for i := range nextDirections {
		nextDirections[i] = newDirection(directionCodes[i], d.oct)
	}
	return nextDirections
}

func CountourTracking(bitmap [][]byte, value byte, oct bool) *ChainCode {
	cc := new(ChainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {

				cc.Start = Point{imageX, imageY}
				cc.Points = []Point{cc.Start}
				prevDirection := newDirection(0, oct)
				currentPoint := NewPoint(cc.Start.X, cc.Start.Y)

				for {
					for _, v := range prevDirection.nextDirections() {
						nextPoint := NewPoint(currentPoint.X+v.d.X, currentPoint.Y+v.d.Y)
						if nextPoint.Y < 0 || nextPoint.X < 0 || len(bitmap) <= nextPoint.Y || len(bitmap[0]) <= nextPoint.X {
							continue
						}
						if bitmap[nextPoint.Y][nextPoint.X] == value {
							cc.Code = append(cc.Code, v.code)
							cc.Points = append(cc.Points, *nextPoint)
							prevDirection = v
							currentPoint = nextPoint
							break
						}
					}
					if cc.Start == *currentPoint {
						break
					}
				}

				return cc
			}
		}
	}
	return cc
}
