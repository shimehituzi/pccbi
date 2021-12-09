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
