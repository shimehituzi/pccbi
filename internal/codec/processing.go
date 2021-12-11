package codec

import (
	"math"
)

// ラベリング
func newLabel(bm bitmap) (label, int) {
	label := make(label, len(bm))
	for y := range label {
		label[y] = make([]int, len(bm[0]))
	}

	lookupTable := []int{}
	counter := 0

	mask := [4]point{
		{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	}

	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				// 周りの 0 以外の値を取得
				al := []int{}
				for _, m := range mask {
					p := point{x + m.x, y + m.y}
					if validPointInt(p, label) && label[p.y][p.x] != 0 {
						al = append(al, label[p.y][p.x])
					}
				}
				if len(al) == 0 {
					// 新規ラベル作成
					counter++
					lookupTable = append(lookupTable, counter)
					label[y][x] = counter
				} else {
					// ルックアップテーブルをたどり一番小さい値を取得する
					root := []int{}
					for _, v := range al {
						k := lookupTable[v-1] - 1
						for k != lookupTable[k]-1 {
							k = lookupTable[k] - 1
						}
						root = append(root, k+1)
					}
					// 最小値を算出
					min := math.MaxInt32
					for _, v := range root {
						if min > lookupTable[v-1] {
							min = lookupTable[v-1]
						}
					}
					// ルックアップテーブル更新
					for _, v := range root {
						lookupTable[v-1] = min
					}
					// ラベリング
					label[y][x] = min
				}
			}
		}
	}

	// ルックアップテーブルを更新
	for i := range lookupTable {
		if i == lookupTable[i]-1 {
			continue
		}
		k := lookupTable[i] - 1
		for k != lookupTable[k]-1 {
			k = lookupTable[k] - 1
		}
		lookupTable[i] = k + 1
	}

	// ルックアップテーブルの値を詰める
	updateTable := []int{}
	for _, v := range lookupTable {
		flag := true
		for _, u := range updateTable {
			if u == v {
				flag = false
			}
		}
		if flag {
			updateTable = append(updateTable, v)
		}
	}
	numLabel := len(updateTable)

	for i, u := range updateTable {
		for j, v := range lookupTable {
			if v == u {
				lookupTable[j] = i + 1
			}
		}
	}

	// label を更新
	for y := range label {
		for x, v := range label[y] {
			if v != 0 {
				newLabel := lookupTable[v-1]
				label[y][x] = newLabel
			}
		}
	}

	return label, numLabel
}

// 輪郭追跡
type orientedPoint struct {
	p    point
	code byte
}

func (p orientedPoint) candidatePoints() [8]orientedPoint {
	firstDirection := byte(p.code + 5)
	nextP := [8]orientedPoint{}
	for i := range nextP {
		d := newDirection((firstDirection + byte(i)) % 8)
		nextP[i].code = d.code
		nextP[i].p = point{p.p.x + d.d.x, p.p.y + d.d.y}
	}
	return nextP
}

func newChaincode(img bitmap, start point, value byte, inner bool) *chaincode {
	current := orientedPoint{start, newDirection(0).code}

	checkP := point{start.x - 1, start.y + 1}
	checkFlag := false
	if !checkP.checkValue(img, value) {
		checkP = start
	}

	code := []byte{}
	points := []point{start}
	for {
		ps := current.candidatePoints()
		for i, candidate := range ps {
			if candidate.p.checkValue(img, value) {
				if inner && (candidate.code%2 == 1) {
					prev := ps[(i-1)%8]
					next := ps[(i+1)%8]
					if prev.p.checkValue(img, 1) && next.p.checkValue(img, 1) {
						continue
					}
				}
				code = append(code, candidate.code)
				points = append(points, candidate.p)
				current = candidate
				if checkP == current.p {
					checkFlag = true
				}
				break
			}
			if i == len(ps)-1 {
				checkFlag = true
			}
		}
		if checkFlag && start == current.p {
			break
		}
	}

	return &chaincode{
		Start:  start,
		Code:   code,
		Points: points,
	}
}

func (p point) checkValue(img bitmap, value byte) bool {
	return validPointByte(p, img) && img[p.y][p.x] == value
}

func newChaincode0(img bitmap, start point, value byte, inner bool) *chaincode {
	currentD := newDirection(0)
	currentP := start

	checkP := point{start.x - 1, start.y + 1}
	if !(validPointByte(checkP, img) && img[checkP.y][checkP.x] == value) {
		checkP = start
	}

	code := []byte{}
	points := []point{start}
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
				code = append(code, nextD.code)
				points = append(points, nextP)
				currentD = nextD
				currentP = nextP
				break
			}
		}
		if start == currentP && checkP.in(points) {
			break
		}
	}

	return &chaincode{
		Start:  start,
		Code:   code,
		Points: points,
	}
}

func newDirection(code byte) direction {
	directions := [8]point{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
	d := directions[code]
	return direction{d, code}
}

func (d direction) nextDirections() [8]direction {
	firstDirection := byte(d.code + 5)
	nextDirections := [8]direction{}
	for i := range nextDirections {
		directionCode := (firstDirection + byte(i)) % 8
		nextDirections[i] = newDirection(directionCode)
	}
	return nextDirections
}

func getChainCodePoints(start point, code []byte) []point {
	d := newDirection(0).d
	p := start
	points := []point{p}

	for i := range code {
		d = newDirection(code[i]).d
		p = point{p.x + d.x, p.y + d.y}
		points = append(points, p)
	}

	return points
}
