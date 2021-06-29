package refactoring

type labeledBitMap [][]int
type labeledBitMaps []labeledBitMap

func newLabeledBitMap(img [][]byte) labeledBitMap {
	ccs := newChainCodes(img)

	lbm := make(labeledBitMap, len(img))
	for i := range lbm {
		lbm[i] = make([]int, len(img[i]))
	}

	label := 0
	for i, cc := range ccs {
		adjacentPoint := getAdjacentPoint(cc, ccs[:i])
		if adjacentPoint != nil {
			for _, point := range cc.points {
				lbm[point.y][point.x] = lbm[adjacentPoint.y][adjacentPoint.x]
			}
		} else {
			label++
			for _, point := range cc.points {
				lbm[point.y][point.x] = label
			}
		}
	}

	return lbm
}

func getAdjacentPoint(cc chainCode, ccs []chainCode) *point {
	if len(ccs) == 0 {
		return nil
	}
	for _, p := range cc.points {
		directions := [8]direction{}
		for i := range directions {
			directions[i] = newDirection(byte(i), true)
		}
		for _, d := range directions {
			adjacnetPoint := point{
				p.x + d.d.x,
				p.y + d.d.y,
			}
			for _, chaincode := range ccs {
				for _, contourPoint := range chaincode.points {
					if adjacnetPoint == contourPoint {
						return &contourPoint
					}
				}
			}
		}
	}
	return nil
}

func newChainCodes(img [][]byte) []chainCode {
	tmp := make([][]byte, len(img))
	for i := range tmp {
		tmp[i] = make([]byte, len(img[i]))
		copy(tmp, img)
	}

	ccs := []chainCode{}
	for {
		cc := contourTracking(tmp, 1, true)
		ccs = append(ccs, cc)
		for _, point := range cc.points {
			tmp[point.y][point.x] = 0
		}
		if isNotExistPoint(tmp) {
			break
		}
	}
	return ccs
}

func isNotExistPoint(img [][]byte) bool {
	for y := range img {
		for x := range img[y] {
			if img[y][x] == 1 {
				return false
			}
		}
	}
	return true
}
