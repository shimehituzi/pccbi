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
	Dx, Dy int
}

func getDirection8() [8]Direction {
	return [8]Direction{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
}

func CountourTracking8(bitmap [][]byte, value byte) *ChainCode {
	cc := new(ChainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {
				cc.Start = Point{imageX, imageY}

				// ================ここが輪郭追跡================
				direction := getDirection8()
				currentPoint := NewPoint(cc.Start.X, cc.Start.Y)
				prevDirection := direction[0]
				cc.Points = []Point{*currentPoint}
				for {
					for _, v := range prevDirection.nextDirection8() {
						nextPoint := NewPoint(currentPoint.X+direction[v].Dx, currentPoint.Y+direction[v].Dy)
						if nextPoint.Y < 0 || nextPoint.X < 0 || len(bitmap) <= nextPoint.Y || len(bitmap[0]) <= nextPoint.X {
							continue
						}
						if bitmap[nextPoint.Y][nextPoint.X] == value {
							currentPoint = nextPoint
							prevDirection = direction[v]
							cc.Points = append(cc.Points, *currentPoint)
							cc.Code = append(cc.Code, prevDirection.toCode8())
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

func (d Direction) toCode8() byte {
	direction := getDirection8()
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

func (d Direction) nextDirection8() [8]int {
	direction := getDirection8()
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

func getDirection4() [4]Direction {
	return [4]Direction{
		{1, 0}, {0, 1}, {-1, 0}, {0, -1},
	}
}

func CountourTracking4(bitmap [][]byte, value byte) *ChainCode {
	cc := new(ChainCode)
	for imageY := range bitmap {
		for imageX := range bitmap[imageY] {
			if bitmap[imageY][imageX] == value {
				cc.Start = Point{imageX, imageY}

				// ================ここが輪郭追跡================
				direction := getDirection4()
				currentPoint := NewPoint(cc.Start.X, cc.Start.Y)
				prevDirection := direction[0]
				cc.Points = []Point{*currentPoint}
				for {
					for _, v := range prevDirection.nextDirection4() {
						nextPoint := NewPoint(currentPoint.X+direction[v].Dx, currentPoint.Y+direction[v].Dy)
						if nextPoint.Y < 0 || nextPoint.X < 0 || len(bitmap) <= nextPoint.Y || len(bitmap[0]) <= nextPoint.X {
							continue
						}
						if bitmap[nextPoint.Y][nextPoint.X] == value {
							currentPoint = nextPoint
							prevDirection = direction[v]
							cc.Points = append(cc.Points, *currentPoint)
							cc.Code = append(cc.Code, prevDirection.toCode4())
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

func (d Direction) toCode4() byte {
	direction := getDirection4()
	switch d {
	case direction[0]:
		return 0
	case direction[1]:
		return 1
	case direction[2]:
		return 2
	case direction[3]:
		return 3
	default:
		panic("その direction は存在しません")
	}
}

func (d Direction) nextDirection4() [4]int {
	direction := getDirection4()
	switch d {
	case direction[0]:
		return [4]int{3, 0, 1, 2}
	case direction[1]:
		return [4]int{0, 1, 2, 3}
	case direction[2]:
		return [4]int{1, 2, 3, 0}
	case direction[3]:
		return [4]int{2, 3, 0, 1}
	default:
		panic("その direction は存在しません")
	}
}
