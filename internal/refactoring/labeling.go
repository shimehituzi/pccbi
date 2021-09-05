package refactoring

import (
	"math"
)

type label intmap

func NewLabel(bm bitmap) label {
	label := make(label, len(bm))
	for y := range label {
		label[y] = make([]int, len(bm[0]))
	}

	lookupTable := []int{}
	counter := 0

	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				al := arroundLavel(x, y, label)
				if len(al) == 0 {
					// 新規ラベル作成
					counter++
					lookupTable = append(lookupTable, counter)
					label[y][x] = counter
				} else {
					// 最小値を算出
					min := math.MaxInt32
					for _, v := range al {
						if min > v {
							min = v
						}
					}
					// ルックアップテーブル更新
					for _, v := range al {
						lookupTable[v-1] = min
					}
					// ラベリング
					label[y][x] = min
				}
			}
		}
	}

	for y := range label {
		for x, v := range label[y] {
			if v != 0 {
				newLabel := lookupTable[v-1]
				label[y][x] = newLabel
			}
		}
	}

	return label
}

func arroundLavel(x, y int, label label) []int {
	mask := []point{
		{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	}

	al := []int{}
	for _, m := range mask {
		p := point{x + m.x, y + m.y}
		if validPoint(p, intmap(label)) && label[p.y][p.x] != 0 {
			al = append(al, label[p.y][p.x])
		}
	}

	return al
}
