package labeling

type Point struct {
	X, Y int
}

func NewPoint(X, Y int) *Point {
	return &Point{X, Y}
}

type Direction struct {
	Dx, Dy int
}

func getDirection() [8]Direction {
	return [8]Direction{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
}

type ChainCode struct {
	Start  Point
	Code   []byte
	Points []Point
}

func CountourTracking(bitmap [][]byte, value byte) *ChainCode {
	cc := new(ChainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {
				cc.Start = Point{imageX, imageY}

				// ================ここが輪郭追跡================
				direction := getDirection()
				currentPoint := NewPoint(cc.Start.X, cc.Start.Y)
				prevDirection := direction[0]
				cc.Points = []Point{*currentPoint}
				for {
					for _, v := range prevDirection.nextDirection() {
						nextPoint := NewPoint(currentPoint.X+direction[v].Dx, currentPoint.Y+direction[v].Dy)
						if nextPoint.Y < 0 || nextPoint.X < 0 || len(bitmap) <= nextPoint.Y || len(bitmap[0]) <= nextPoint.X {
							continue
						}
						if bitmap[nextPoint.Y][nextPoint.X] == value {
							currentPoint = nextPoint
							prevDirection = direction[v]
							cc.Points = append(cc.Points, *currentPoint)
							cc.Code = append(cc.Code, prevDirection.toCode())
							break
						}
					}
					if cc.Start == *currentPoint {
						break
					}
				}
				// ==============================================

				return cc
			}
		}
	}
	return cc
}

func (d Direction) toCode() byte {
	direction := getDirection()
	switch d {
	case direction[0]:
		return 0
	case direction[1]:
		return 1
	case direction[2]:
		return 2
	case direction[3]:
		return 3
	case direction[4]:
		return 4
	case direction[5]:
		return 5
	case direction[6]:
		return 6
	case direction[7]:
		return 7
	default:
		panic("その direction は存在しません")
	}
}

func (d Direction) nextDirection() [8]int {
	direction := getDirection()
	switch d {
	case direction[0]:
		return [8]int{5, 6, 7, 0, 1, 2, 3, 4}
	case direction[1]:
		return [8]int{6, 7, 0, 1, 2, 3, 4, 5}
	case direction[2]:
		return [8]int{7, 0, 1, 2, 3, 4, 5, 6}
	case direction[3]:
		return [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	case direction[4]:
		return [8]int{1, 2, 3, 4, 5, 6, 7, 0}
	case direction[5]:
		return [8]int{2, 3, 4, 5, 6, 7, 0, 1}
	case direction[6]:
		return [8]int{3, 4, 5, 6, 7, 0, 1, 2}
	case direction[7]:
		return [8]int{4, 5, 6, 7, 0, 1, 2, 3}
	default:
		panic("その direction は存在しません")
	}
}
