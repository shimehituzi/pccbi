package labeling

import "fmt"

type Direction struct {
	Dx, Dy int
}

func returnDirection() [8]Direction {
	return [8]Direction{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
}

func (d Direction) toCode() byte {
	direction := returnDirection()
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
	direction := returnDirection()
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

func CountourTracking(lbm *LabeledBitMap) *ChainCode {
	cc := new(ChainCode)
	for imageY := range lbm.Image {
		for imageX := range lbm.Image[imageY] {
			if lbm.Image[imageY][imageX] == 1 {
				cc.Start = Point{imageX, imageY}
				InnerCountourTracking(cc, lbm)
				break
			}
		}
		break
	}
	fmt.Println("cc.start", cc.Start)
	fmt.Println("cc.Code", cc.Code)
	fmt.Println("cc.Points", cc.Points)
	return cc
}

func InnerCountourTracking(cc *ChainCode, lbm *LabeledBitMap) {
	direction := returnDirection()
	currentPoint := NewPoint(cc.Start.X, cc.Start.Y)
	prevDirection := direction[0]
	cc.Points = []Point{*currentPoint}
	// cc.Code = []byte{prevDirection.toCode()}
	for {
		for _, v := range prevDirection.nextDirection() {
			nextPoint := NewPoint(currentPoint.X+direction[v].Dx, currentPoint.Y+direction[v].Dy)
			if nextPoint.Y < 0 || nextPoint.X < 0 || len(lbm.Image) <= nextPoint.Y || len(lbm.Image[0]) <= nextPoint.X {
				continue
			}
			if lbm.Image[nextPoint.Y][nextPoint.X] == 1 {
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
}
